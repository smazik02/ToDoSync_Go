package server

import (
	"bufio"
	"bytes"
	"context"
	"database/sql"
	"errors"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"todosync_go/internal/database"
	"todosync_go/internal/parser"
	"todosync_go/internal/services"
	"todosync_go/internal/shared"

	"golang.org/x/sync/errgroup"
)

const BUFSIZE = 4096
const MAX_BUFFER_SIZE = 8 * 1024 * 1024

type Server struct {
	sync.RWMutex
	eg             *errgroup.Group
	listener       *net.TCPListener
	clients        map[net.Conn]*shared.Client
	connections    chan net.Conn
	serviceGateway *services.ServiceGateway
}

func NewServer(port int, serviceGateway *services.ServiceGateway, db *sql.DB) (*Server, error) {
	addr := net.TCPAddr{Port: port}

	ln, err := net.ListenTCP("tcp", &addr)
	if err != nil {
		return nil, err
	}

	database.CreateTables(db)

	return &Server{
		listener:       ln,
		clients:        make(map[net.Conn]*shared.Client),
		connections:    make(chan net.Conn),
		serviceGateway: serviceGateway,
	}, nil
}

func (s *Server) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	eg, egCtx := errgroup.WithContext(ctx)
	s.eg = eg

	s.eg.Go(func() error {
		<-egCtx.Done()
		return s.listener.Close()
	})

	s.eg.Go(func() error {
		return s.acceptConnections(egCtx)
	})

	s.eg.Go(func() error {
		s.handleTerminal(cancel)
		return nil
	})

	for {
		select {
		case incoming := <-s.connections:
			s.eg.Go(func() error {
				s.handleConnection(ctx, incoming)
				return nil
			})
		case <-egCtx.Done():
			return s.eg.Wait()
		}
	}
}

func (s *Server) handleTerminal(cancel context.CancelFunc) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		text := scanner.Text()

		if err := scanner.Err(); err != nil {
			continue
		}

		if command := strings.ToLower(strings.TrimRight(text, "\r\n")); command == "q" {
			cancel()
			return
		}
	}
}

func (s *Server) acceptConnections(ctx context.Context) error {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return err
		}

		select {
		case s.connections <- conn:
		case <-ctx.Done():
			conn.Close()
			return nil
		}
	}
}

func (s *Server) handleConnection(ctx context.Context, connection net.Conn) {
	defer connection.Close()

	clientAddress := connection.RemoteAddr().String()
	log.Printf("[%s] New connection\n", clientAddress)

	client := &shared.Client{Connection: connection, UserId: -1}
	s.Lock()
	s.clients[connection] = client
	s.Unlock()

	done := make(chan any)
	defer close(done)

	defer func() {
		s.Lock()
		delete(s.clients, connection)
		s.Unlock()
	}()

	buf := make([]byte, BUFSIZE)

	conChan, errChan := make(chan []byte, 1), make(chan error, 1)

	s.eg.Go(func() error {
		for {
			cnt, err := connection.Read(buf)
			if err != nil {
				select {
				case errChan <- err:
				case <-done:
				}
				return nil
			}

			data := make([]byte, cnt)
			copy(data, buf[:cnt])

			select {
			case conChan <- data:
			case <-done:
				return nil
			}
		}
	})

	for {
		select {
		case <-ctx.Done():
			connection.Write([]byte("Disconnecting!"))
			return

		case err := <-errChan:
			if errors.Is(err, io.EOF) {
				connection.Write([]byte("Disconnecting!"))
				log.Printf("[%s] Disconnected\n", clientAddress)
			} else {
				log.Printf("[%s] Error: %s\n", clientAddress, err.Error())
			}
			return

		case data := <-conChan:
			if client.Buffer.Len()+len(data) > MAX_BUFFER_SIZE {
				log.Printf("[%s] Client exceeded max buffer size. Disconnecting.\n", clientAddress)
				connection.Write([]byte("Disconnecting!"))
				return
			}

			client.Buffer.Write(data)
			separator := []byte("\n\n")

			for {
				unreadBytes := client.Buffer.Bytes()

				idx := bytes.Index(unreadBytes, separator)
				if idx == -1 {
					break
				}

				messageStr := string(unreadBytes[:idx])

				client.Buffer.Next(idx + len(separator))

				parsedMessage, err := parser.ProcessRequest(messageStr)
				if err != nil {
					log.Printf("[%s] Parser error occured\n", clientAddress)
					connection.Write([]byte(err.Error()))
					log.Printf("[%s] Error message sent to client", clientAddress)
					continue
				}

				// log.Printf("[%s] Parsed message: %s|%s\n", clientAddress, parsedMessage.ResourceMethod, string(parsedMessage.Payload))
				response, err := s.serviceGateway.Direct(parsedMessage, client)
				if err != nil {
					log.Printf("[%s] Service error occured\n", clientAddress)
					connection.Write([]byte(err.Error()))
					log.Printf("[%s] Error message sent to client", clientAddress)
					continue
				}

				log.Printf("[%s] Sent response to client\n", clientAddress)
				connection.Write(response.Message)
			}
		}
	}
}
