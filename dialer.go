package aws

import (
	"net"
	"os"
	"sync"
	//  "log"
)

// Dev notes: lower-case (private) functions assume the lock is held,
// upper-case functions should use a defer lock.Unlock to ensure
// underlying dialer/socket panics will not leave locks hanging.

var ErrUnderlyingNotconnected = os.NewError("Underlying socket is not connected")

// A Dialer is usually a closuer that
// is pre-configured to the callers tastes.
//
// (see URLDialer for an example/default generator)
type Dialer func() (net.Conn, os.Error)

// A Reusable conn is a syncronized structure around a
// Dialer / net.Conn pair.  All net.Conn calls are wrapped
// around the underlying structure.  Errors are bubbled
// up, and trigger closure of the underlying socket (to
// be reopened on the next call)
type ReusableConn struct {
	lock         *sync.Mutex
	dialer       Dialer
	conn         net.Conn
	readTimeout  int64
	writeTimeout int64
}

const (
	_UNSET_TIMEOUT int64 = -1
)

// Create a new reusable connection with a sepcific dialer.
func NewReusableConnection(d Dialer) (c *ReusableConn) {
	return &ReusableConn{
		dialer:       d,
		conn:         nil,
		lock:         &sync.Mutex{},
		readTimeout:  _UNSET_TIMEOUT,
		writeTimeout: _UNSET_TIMEOUT,
	}
}

// Dial is idempotent, and safe to call;
func (self *ReusableConn) Dial() (err os.Error) {
	// log.Printf("Public Dial() called")
	self.lock.Lock()
	defer self.lock.Unlock()
	return self.dial()
}

// Dial will redial if conn is nil, and set
// timeouts if they've been set by the caller.
// 
// It simply returns nil if the socket appears already connected
func (self *ReusableConn) dial() (err os.Error) {
	// log.Printf("Private dial() called (%v)", self.conn)
	if self.conn == nil {
		self.conn, err = self.dialer()
		if err == nil && self.readTimeout != _UNSET_TIMEOUT {
			err = self.setReadTimeout(self.readTimeout)
		}
		if err == nil && self.writeTimeout != _UNSET_TIMEOUT {
			err = self.setWriteTimeout(self.writeTimeout)
		}
	}
	// log.Printf("Private dial() complete (%v)", self.conn)
	return
}

func (self *ReusableConn) close() (err os.Error) {
	if self.conn != nil {
		err = self.conn.Close()
		self.conn = nil
	}
	return
}

// Unlike close on a traditional socket, no error
// is raised if you close a closed (nil) connection.
func (self *ReusableConn) Close() (err os.Error) {
	self.lock.Lock()
	defer self.lock.Unlock()
	return self.close()
}

// TODO: What's an appropriate responsde when we're not connected?
// ATM, we return whatever the other side says, or the nil net.Addr.
func (self *ReusableConn) RemoteAddr() (a net.Addr) {
	self.lock.Lock()
	defer self.lock.Unlock()
	if self.conn != nil {
		a = self.conn.RemoteAddr()
	}
	return
}

// See RemoteAddr for notes.
func (self *ReusableConn) LocalAddr() (a net.Addr) {
	self.lock.Lock()
	defer self.lock.Unlock()
	if self.conn != nil {
		a = self.conn.RemoteAddr()
	}
	return
}

func (self *ReusableConn) read(in []byte) (n int, err os.Error) {
	err = self.dial()
	if err == nil {
		n, err = self.conn.Read(in)
		if err != nil {
			self.close()
		}
	}
	return
}

func (self *ReusableConn) write(in []byte) (n int, err os.Error) {
	err = self.dial()
	if err == nil {
		n, err = self.conn.Write(in)
		if err != nil {
			self.close()
		}
	}
	return
}


// Read from the underlying connection, triggering a dial if needed.
// NB: For the expected case (HTTP), this shouldn't happen before the
// first Write.
func (self *ReusableConn) Read(in []byte) (n int, err os.Error) {
	self.lock.Lock()
	defer self.lock.Unlock()
	return self.read(in)
}

// Write to the underlying connection, triggering a dial if needed.
func (self *ReusableConn) Write(out []byte) (n int, err os.Error) {
	self.lock.Lock()
	defer self.lock.Unlock()
	return self.write(out)
}


func (self *ReusableConn) setReadTimeout(t int64) (err os.Error) {
	err = self.dial()
	if err == nil {
		err = self.conn.SetReadTimeout(t)
		if err == nil {
			self.readTimeout = t
		}
	}
	return
}

func (self *ReusableConn) setWriteTimeout(t int64) (err os.Error) {
	err = self.dial()
	if err == nil {
		err = self.conn.SetWriteTimeout(t)
		if err == nil {
			self.writeTimeout = t
		}
	}
	return
}


// Sets the read timeout on the underlying socket, as well
// as an internal flag for any future re-opened connections.
func (self *ReusableConn) SetReadTimeout(t int64) (err os.Error) {
	self.lock.Lock()
	defer self.lock.Unlock()
	return self.setReadTimeout(t)
}

// Sets the write timeout on the underlying socket, as well
// as an internal flag for any future re-opened connections.
func (self *ReusableConn) SetWriteTimeout(t int64) (err os.Error) {
	self.lock.Lock()
	defer self.lock.Unlock()
	return self.setWriteTimeout(t)
}

// Conveinience function for Set(read|write)timeout
func (self *ReusableConn) SetTimeout(t int64) (err os.Error) {
	err = self.SetReadTimeout(t)
	if err == nil {
		err = self.SetWriteTimeout(t)
	}
	return
}
