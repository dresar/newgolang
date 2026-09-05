package services

import (
	"context"
	"errors"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

type WhatsAppStatus struct {
	Connected   bool   `json:"connected"`
	LastEvent   string `json:"last_event"`
	LastQRCode  string `json:"last_qr_code"`
	PhoneNumber string `json:"phone_number"`
}

type WhatsAppService struct {
	container   *sqlstore.Container
	client      *whatsmeow.Client
	log         waLog.Logger
	statusMu    sync.RWMutex
	connected   bool
	lastEvent   string
	lastQRCode  string
	phoneNumber string
	startOnce   sync.Once
}

var waServiceInstance *WhatsAppService
var waServiceOnce sync.Once

func NewWhatsAppService(dbPath string) (*WhatsAppService, error) {
	var err error
	waServiceOnce.Do(func() {
		log := waLog.Stdout("WhatsApp", "INFO", true)
		container, cErr := sqlstore.New(context.Background(), "sqlite3", "file:"+dbPath+"?_foreign_keys=on", log)
		if cErr != nil {
			err = cErr
			return
		}
		device, dErr := container.GetFirstDevice(context.Background())
		if dErr != nil {
			err = dErr
			return
		}
		client := whatsmeow.NewClient(device, log)
		svc := &WhatsAppService{
			container: container,
			client:    client,
			log:       log,
		}
		client.AddEventHandler(svc.handleEvent)
		waServiceInstance = svc
	})
	if err != nil {
		return nil, err
	}
	return waServiceInstance, nil
}

func (s *WhatsAppService) handleEvent(evt interface{}) {
	switch evt.(type) {
	case *events.Connected:
		s.statusMu.Lock()
		s.connected = true
		s.lastEvent = "connected"
		if s.client.Store.ID != nil {
			s.phoneNumber = s.client.Store.ID.User
		}
		s.statusMu.Unlock()
	case *events.Disconnected:
		s.statusMu.Lock()
		s.connected = false
		s.lastEvent = "disconnected"
		s.statusMu.Unlock()
	case *events.StreamReplaced:
		s.statusMu.Lock()
		s.connected = false
		s.lastEvent = "replaced"
		s.statusMu.Unlock()
	}
}

func (s *WhatsAppService) Start() error {
	var err error
	s.startOnce.Do(func() {
		if s.client.Store.ID == nil {
			err = s.startWithQRLogin()
		} else {
			err = s.client.Connect()
			if err == nil {
				s.statusMu.Lock()
				s.lastEvent = "connected_existing"
				s.connected = true
				s.statusMu.Unlock()
			}
		}
	})
	return err
}

func (s *WhatsAppService) startWithQRLogin() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	qrChan, err := s.client.GetQRChannel(ctx)
	if err != nil {
		return err
	}
	go func() {
		for evt := range qrChan {
			if evt.Event == "code" {
				s.statusMu.Lock()
				s.lastQRCode = evt.Code
				s.lastEvent = "qr_code"
				s.statusMu.Unlock()
			} else {
				s.statusMu.Lock()
				s.lastEvent = evt.Event
				if evt.Event == "timeout" || evt.Event == "error" {
					s.lastQRCode = ""
				}
				s.statusMu.Unlock()
			}
		}
	}()
	return s.client.Connect()
}

func (s *WhatsAppService) Reconnect() error {
	if s.client == nil {
		return errors.New("client not initialized")
	}
	s.client.Disconnect()
	time.Sleep(time.Second)
	err := s.client.Connect()
	if err != nil {
		return err
	}
	s.statusMu.Lock()
	s.connected = true
	s.lastEvent = "reconnected"
	s.statusMu.Unlock()
	return nil
}

func (s *WhatsAppService) Status() WhatsAppStatus {
	s.statusMu.RLock()
	defer s.statusMu.RUnlock()
	return WhatsAppStatus{
		Connected:   s.connected,
		LastEvent:   s.lastEvent,
		LastQRCode:  s.lastQRCode,
		PhoneNumber: s.phoneNumber,
	}
}

func (s *WhatsAppService) SendTextMessage(to string, message string) error {
	if s.client == nil {
		return errors.New("client not initialized")
	}
	jid := types.JID{
		User:   to,
		Server: types.DefaultUserServer,
	}
	msg := &waE2E.Message{
		Conversation: proto.String(message),
	}
	_, err := s.client.SendMessage(context.Background(), jid, msg, whatsmeow.SendRequestExtra{})
	return err
}
