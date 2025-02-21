package gomail

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"net"
	"net/smtp"
	"reflect"
	"testing"
)

const (
	testPort    = 587
	testSSLPort = 465
)

var (
	testConn    = &net.TCPConn{}
	testTLSConn = tls.Client(testConn, &tls.Config{InsecureSkipVerify: true})
	testConfig  = &tls.Config{InsecureSkipVerify: true}
	testAuth    = smtp.PlainAuth("", testUser, testPwd, testHost)
)

func TestDialer(t *testing.T) {
	d := NewDialer(testHost, testPort, "user", "pwd")
	testSendMail(t, d, []string{
		"Extension STARTTLS",
		"StartTLS",
		"Extension AUTH",
		"Auth",
		"Mail " + testFrom,
		"Rcpt " + testTo1,
		"Rcpt " + testTo2,
		"Data",
		"Write message",
		"Close writer",
		"Quit",
		"Close",
	})
}

func TestDialerSSL(t *testing.T) {
	d := NewDialer(testHost, testSSLPort, "user", "pwd")
	testSendMail(t, d, []string{
		"Extension AUTH",
		"Auth",
		"Mail " + testFrom,
		"Rcpt " + testTo1,
		"Rcpt " + testTo2,
		"Data",
		"Write message",
		"Close writer",
		"Quit",
		"Close",
	})
}

func TestDialerConfig(t *testing.T) {
	d := NewDialer(testHost, testPort, "user", "pwd")
	d.LocalName = "test"
	d.TLSConfig = testConfig
	testSendMail(t, d, []string{
		"Hello test",
		"Extension STARTTLS",
		"StartTLS",
		"Extension AUTH",
		"Auth",
		"Mail " + testFrom,
		"Rcpt " + testTo1,
		"Rcpt " + testTo2,
		"Data",
		"Write message",
		"Close writer",
		"Quit",
		"Close",
	})
}

func TestDialerSSLConfig(t *testing.T) {
	d := NewDialer(testHost, testSSLPort, "user", "pwd")
	d.LocalName = "test"
	d.TLSConfig = testConfig
	testSendMail(t, d, []string{
		"Hello test",
		"Extension AUTH",
		"Auth",
		"Mail " + testFrom,
		"Rcpt " + testTo1,
		"Rcpt " + testTo2,
		"Data",
		"Write message",
		"Close writer",
		"Quit",
		"Close",
	})
}

func TestDialerNoStartTLS(t *testing.T) {
	d := NewDialer(testHost, testPort, "user", "pwd")
	d.StartTLSPolicy = NoStartTLS
	testSendMail(t, d, []string{
		"Extension AUTH",
		"Auth",
		"Mail " + testFrom,
		"Rcpt " + testTo1,
		"Rcpt " + testTo2,
		"Data",
		"Write message",
		"Close writer",
		"Quit",
		"Close",
	})
}

func TestDialerOpportunisticStartTLS(t *testing.T) {
	d := NewDialer(testHost, testPort, "user", "pwd")
	d.StartTLSPolicy = OpportunisticStartTLS
	testSendMail(t, d, []string{
		"Extension STARTTLS",
		"StartTLS",
		"Extension AUTH",
		"Auth",
		"Mail " + testFrom,
		"Rcpt " + testTo1,
		"Rcpt " + testTo2,
		"Data",
		"Write message",
		"Close writer",
		"Quit",
		"Close",
	})

	if OpportunisticStartTLS != 0 {
		t.Errorf("OpportunisticStartTLS: expected 0, got %d",
			OpportunisticStartTLS)
	}
}

func TestDialerOpportunisticStartTLSUnsupported(t *testing.T) {
	d := NewDialer(testHost, testPort, "user", "pwd")
	d.StartTLSPolicy = OpportunisticStartTLS
	testSendMailStartTLSUnsupported(t, d, []string{
		"Extension STARTTLS",
		"Extension AUTH",
		"Auth",
		"Mail " + testFrom,
		"Rcpt " + testTo1,
		"Rcpt " + testTo2,
		"Data",
		"Write message",
		"Close writer",
		"Quit",
		"Close",
	})
}

func TestDialerMandatoryStartTLS(t *testing.T) {
	d := NewDialer(testHost, testPort, "user", "pwd")
	d.StartTLSPolicy = MandatoryStartTLS
	testSendMail(t, d, []string{
		"Extension STARTTLS",
		"StartTLS",
		"Extension AUTH",
		"Auth",
		"Mail " + testFrom,
		"Rcpt " + testTo1,
		"Rcpt " + testTo2,
		"Data",
		"Write message",
		"Close writer",
		"Quit",
		"Close",
	})
}

func TestDialerMandatoryStartTLSUnsupported(t *testing.T) {
	d := NewDialer(testHost, testPort, "user", "pwd")
	d.StartTLSPolicy = MandatoryStartTLS

	testClient := &mockClient{
		t:        t,
		addr:     addr(d.Host, d.Port),
		config:   d.TLSConfig,
		startTLS: false,
		timeout:  true,
	}

	err := doTestSendMail(t, d, testClient, []string{
		"Extension STARTTLS",
	})

	if _, ok := err.(StartTLSUnsupportedError); !ok {
		t.Errorf("expected StartTLSUnsupportedError, but got: %s",
			reflect.TypeOf(err).Name())
	}

	expected := "gomail: MandatoryStartTLS required, " +
		"but SMTP server does not support STARTTLS"
	if err.Error() != expected {
		t.Errorf("expected %s, but got: %s", expected, err)
	}
}

func TestDialerNoAuth(t *testing.T) {
	d := &Dialer{
		Host: testHost,
		Port: testPort,
	}
	testSendMail(t, d, []string{
		"Extension STARTTLS",
		"StartTLS",
		"Mail " + testFrom,
		"Rcpt " + testTo1,
		"Rcpt " + testTo2,
		"Data",
		"Write message",
		"Close writer",
		"Quit",
		"Close",
	})
}

func TestDialerTimeout(t *testing.T) {
	d := &Dialer{
		Host:         testHost,
		Port:         testPort,
		RetryFailure: true,
	}
	testSendMailTimeout(t, d, []string{
		"Extension STARTTLS",
		"StartTLS",
		"Mail " + testFrom,
		"Extension STARTTLS",
		"StartTLS",
		"Mail " + testFrom,
		"Rcpt " + testTo1,
		"Rcpt " + testTo2,
		"Data",
		"Write message",
		"Close writer",
		"Quit",
		"Close",
	})
}

func TestDialerTimeoutNoRetry(t *testing.T) {
	d := &Dialer{
		Host:         testHost,
		Port:         testPort,
		RetryFailure: false,
	}
	testClient := &mockClient{
		t:        t,
		addr:     addr(d.Host, d.Port),
		config:   d.TLSConfig,
		startTLS: true,
		timeout:  true,
	}

	err := doTestSendMail(t, d, testClient, []string{
		"Extension STARTTLS",
		"StartTLS",
		"Mail " + testFrom,
		"Quit",
	})

	if err.Error() != "gomail: could not send email 1: EOF" {
		t.Error("expected to have got EOF, but got:", err)
	}
}

type mockClient struct {
	t        *testing.T
	i        int
	want     []string
	addr     string
	config   *tls.Config
	startTLS bool
	timeout  bool
}

func (c *mockClient) Hello(localName string) error {
	c.do("Hello " + localName)
	return nil
}

func (c *mockClient) Extension(ext string) (bool, string) {
	c.do("Extension " + ext)
	ok := true
	if ext == "STARTTLS" {
		ok = c.startTLS
	}
	return ok, ""
}

func (c *mockClient) StartTLS(config *tls.Config) error {
	assertConfig(c.t, config, c.config)
	c.do("StartTLS")
	return nil
}

func (c *mockClient) Auth(a smtp.Auth) error {
	if !reflect.DeepEqual(a, testAuth) {
		c.t.Errorf("Invalid auth, got %#v, want %#v", a, testAuth)
	}
	c.do("Auth")
	return nil
}

func (c *mockClient) Mail(from string) error {
	c.do("Mail " + from)
	if c.timeout {
		c.timeout = false
		return io.EOF
	}
	return nil
}

func (c *mockClient) Rcpt(to string) error {
	c.do("Rcpt " + to)
	return nil
}

func (c *mockClient) Data() (io.WriteCloser, error) {
	c.do("Data")
	return &mockWriter{c: c, want: testMsg}, nil
}

func (c *mockClient) Quit() error {
	c.do("Quit")
	return nil
}

func (c *mockClient) Close() error {
	c.do("Close")
	return nil
}

func (c *mockClient) do(cmd string) {
	if c.i >= len(c.want) {
		c.t.Fatalf("Invalid command %q", cmd)
	}

	if cmd != c.want[c.i] {
		c.t.Fatalf("Invalid command, got %q, want %q", cmd, c.want[c.i])
	}
	c.i++
}

type mockWriter struct {
	want string
	c    *mockClient
	buf  bytes.Buffer
}

func (w *mockWriter) Write(p []byte) (int, error) {
	if w.buf.Len() == 0 {
		w.c.do("Write message")
	}
	w.buf.Write(p)
	return len(p), nil
}

func (w *mockWriter) Close() error {
	compareBodies(w.c.t, w.buf.String(), w.want)
	w.c.do("Close writer")
	return nil
}

func testSendMail(t *testing.T, d *Dialer, want []string) {
	testClient := &mockClient{
		t:        t,
		addr:     addr(d.Host, d.Port),
		config:   d.TLSConfig,
		startTLS: true,
		timeout:  false,
	}

	if err := doTestSendMail(t, d, testClient, want); err != nil {
		t.Error(err)
	}
}

func testSendMailStartTLSUnsupported(t *testing.T, d *Dialer, want []string) {
	testClient := &mockClient{
		t:        t,
		addr:     addr(d.Host, d.Port),
		config:   d.TLSConfig,
		startTLS: false,
		timeout:  false,
	}

	if err := doTestSendMail(t, d, testClient, want); err != nil {
		t.Error(err)
	}
}

func testSendMailTimeout(t *testing.T, d *Dialer, want []string) {
	testClient := &mockClient{
		t:        t,
		addr:     addr(d.Host, d.Port),
		config:   d.TLSConfig,
		startTLS: true,
		timeout:  true,
	}

	if err := doTestSendMail(t, d, testClient, want); err != nil {
		t.Error(err)
	}
}

func doTestSendMail(t *testing.T, d *Dialer, testClient *mockClient, want []string) error {
	testClient.want = want

	dialContext = func(ctx context.Context, d *Dialer) (net.Conn, error) {
		return testConn, nil
	}

	tlsClient = func(conn net.Conn, config *tls.Config) *tls.Conn {
		if conn != testConn {
			t.Errorf("Invalid conn, got %#v, want %#v", conn, testConn)
		}
		assertConfig(t, config, testClient.config)
		return testTLSConn
	}

	smtpNewClient = func(conn net.Conn, host string) (smtpClient, error) {
		if host != testHost {
			t.Errorf("Invalid host, got %q, want %q", host, testHost)
		}
		return testClient, nil
	}

	return d.DialAndSend(context.Background(), getTestMessage())
}

func assertConfig(t *testing.T, got, want *tls.Config) {
	if want == nil {
		want = &tls.Config{ServerName: testHost}
	}
	if got.ServerName != want.ServerName {
		t.Errorf("Invalid field ServerName in config, got %q, want %q", got.ServerName, want.ServerName)
	}
	if got.InsecureSkipVerify != want.InsecureSkipVerify {
		t.Errorf("Invalid field InsecureSkipVerify in config, got %v, want %v", got.InsecureSkipVerify, want.InsecureSkipVerify)
	}
}
