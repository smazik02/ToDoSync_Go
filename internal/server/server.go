package server

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"todosync_go/utils"
)

const BUFSIZE = 4096

type Server struct {
	sync.RWMutex
	wg          sync.WaitGroup
	listener    *net.TCPListener
	clients     map[net.Conn]Client
	connections chan net.Conn
	terminate   chan any
	shutdown    chan any
}

type Client struct {
	connection net.Conn
	buffer     strings.Builder
	userId     int
}

func (c Client) isLoggedIn() bool {
	return c.userId != -1
}

func NewServer(port int) (*Server, error) {
	addr := net.TCPAddr{Port: port}

	ln, err := net.ListenTCP("tcp", &addr)
	if err != nil {
		return nil, err
	}

	return &Server{
		wg:          sync.WaitGroup{},
		listener:    ln,
		clients:     make(map[net.Conn]Client),
		connections: make(chan net.Conn),
		terminate:   make(chan any),
		shutdown:    make(chan any),
	}, nil
}

func (s *Server) Run() {
	s.wg.Add(2)
	go s.handleTerminal()
	go s.acceptConnections()

	for {
		select {
		case incoming := <-s.connections:
			s.wg.Add(1)
			go s.handleConnection(incoming)
		case <-s.terminate:
			s.Stop()
			return
		}
	}
}

func (s *Server) Stop() {
	close(s.shutdown)
	s.listener.Close()
	s.wg.Wait()
}

func (s *Server) handleTerminal() {
	defer s.wg.Done()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		text := scanner.Text()

		if err := scanner.Err(); err != nil {
			continue
		}

		if command := strings.ToLower(strings.TrimRight(text, "\r\n")); command == "q" {
			close(s.terminate)
			return
		}
	}
}

func (s *Server) acceptConnections() {
	defer s.wg.Done()

	for {
		select {
		case <-s.shutdown:
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				continue
			}

			s.connections <- conn
		}
	}
}

func (s *Server) handleConnection(connection net.Conn) {
	defer s.wg.Done()
	defer connection.Close()

	clientAddress := connection.RemoteAddr().String()

	log.Printf("[%s] New connection\n", clientAddress)

	client := Client{connection: connection, userId: -1}
	s.Lock()
	s.clients[connection] = client
	s.Unlock()

	buf := make([]byte, BUFSIZE)

	conChan, errChan := make(chan int, 1), make(chan error, 1)
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()

		for {
			cnt, err := connection.Read(buf)
			if err != nil {
				errChan <- err
				return
			}
			conChan <- cnt
		}
	}()

	for {
		select {
		case <-s.shutdown:
			connection.Write([]byte("Disconnecting!"))
			s.Lock()
			delete(s.clients, connection)
			s.Unlock()
			return
		case err := <-errChan:
			if errors.Is(err, io.EOF) {
				connection.Write([]byte("Disconnecting!"))
				log.Printf("[%s] Disconnected\n", clientAddress)
			} else {
				log.Printf("[%s] Error: %s\n", clientAddress, err.Error())
			}
			s.Lock()
			delete(s.clients, connection)
			s.Unlock()
			return
		case cnt := <-conChan:
			client.buffer.Write(buf[:cnt])
			messages := strings.Split(client.buffer.String(), "\n\n")
			client.buffer.Reset()
			client.buffer.WriteString(messages[len(messages)-1])
			messages = messages[:len(messages)-1]

			for _, message := range messages {
				parsedMessage, err := utils.ProcessRequest(message)
				if err != nil {
					log.Printf("[%s] Parser error occured\n", clientAddress)
					connection.Write([]byte(err.Error()))
					log.Printf("[%s] Error message sent to client", clientAddress)
					continue
				}

				log.Printf("[%s] Parsed message: %s|%s\n", clientAddress, parsedMessage.ResourceMethod, string(parsedMessage.Payload))
			}
		}
	}
}
