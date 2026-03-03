package imap

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/icewarp/warpctl/internal/config"
	"github.com/icewarp/warpctl/internal/logger"
)

type Client struct {
	conn   *client.Client
	config *config.Config
}

func NewClient(cfg *config.Config) (*Client, error) {
	addr := cfg.GetServerAddress()
	host := cfg.Server.Host

	tlsConfig := &tls.Config{
		ServerName:         host,
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: !cfg.IMAP.TLSVerify,
	}

	logger.Info("Connecting to IMAP server", "address", addr)

	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: 15 * time.Second}, "tcp", addr, tlsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect via TLS: %w", err)
	}

	c, err := client.New(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to create IMAP client: %w", err)
	}

	return &Client{conn: c, config: cfg}, nil
}

func (c *Client) Login() error {
	logger.Info("Logging in", "username", c.config.Auth.Username)
	if err := c.conn.Login(c.config.Auth.Username, c.config.Auth.Password); err != nil {
		return fmt.Errorf("login failed: %w", err)
	}
	logger.Info("Login successful")
	return nil
}

func (c *Client) Logout() error {
	if c.conn != nil {
		logger.Info("Logging out")
		return c.conn.Logout()
	}
	return nil
}

func (c *Client) SelectMailbox(mailbox string) (*imap.MailboxStatus, error) {
	logger.Info("Selecting mailbox", "mailbox", mailbox)
	mbox, err := c.conn.Select(mailbox, false)
	if err != nil {
		return nil, fmt.Errorf("failed to select %s: %w", mailbox, err)
	}
	return mbox, nil
}

func (c *Client) GetMessageCount(mailbox string) (uint32, error) {
	mbox, err := c.SelectMailbox(mailbox)
	if err != nil {
		return 0, err
	}
	return mbox.Messages, nil
}

func (c *Client) SearchMessages(mailbox string) ([]uint32, error) {
	mbox, err := c.SelectMailbox(mailbox)
	if err != nil {
		return nil, err
	}

	if mbox.Messages == 0 {
		return []uint32{}, nil
	}

	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{}
	criteria.Text = []string{}
	criteria.Header = map[string][]string{}

	ids, err := c.conn.Search(criteria)
	if err != nil {
		return nil, fmt.Errorf("failed to search messages: %w", err)
	}

	return ids, nil
}

func (c *Client) DeleteMessages(mailbox string, ids []uint32) (int, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	deleted := 0
	batchSize := c.config.IMAP.BatchSize
	if batchSize <= 0 {
		batchSize = 5000
	}

	logger.Info("Starting message deletion", "total", len(ids), "batch_size", batchSize)

	for i := 0; i < len(ids); i += batchSize {
		end := i + batchSize
		if end > len(ids) {
			end = len(ids)
		}

		seq := new(imap.SeqSet)
		seq.AddNum(ids[i:end]...)

		item := imap.FormatFlagsOp(imap.AddFlags, true)
		flags := []interface{}{imap.DeletedFlag}

		if err := c.conn.Store(seq, item, flags, nil); err != nil {
			return deleted, fmt.Errorf("failed to mark messages as deleted: %w", err)
		}

		if err := c.conn.Expunge(nil); err != nil {
			return deleted, fmt.Errorf("failed to expunge: %w", err)
		}

		deleted += (end - i)
		logger.Info("Batch deleted", "count", end-i, "total_deleted", deleted)
	}

	return deleted, nil
}
