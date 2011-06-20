package aws

import (
  "net"
  "os"
  "sync"
)

var ErrUnderlyingNotconnected = os.NewError("Underlying socket is not connected")

// A Dialer is usually a closuer that
// is pre-configured to the callers tastes

type Dialer func()(net.Conn, os.Error)

type ReusableConn struct {
  lock *sync.Mutex
  dialer Dialer
  conn net.Conn
  readTimeout int64
  writeTimeout int64
}

const (
  _UNSET_TIMEOUT int64 = -1
)

func NewReusableConnection(d Dialer)(c *ReusableConn){
  return &ReusableConn{
    dialer: d,
    conn: nil,
    lock: &sync.Mutex{},
    readTimeout: _UNSET_TIMEOUT,
    writeTimeout: _UNSET_TIMEOUT,
  }
}

// Dial is idempotent, and safe to call;
func (self *ReusableConn)Dial()(err os.Error){
  self.lock.Lock()
  defer self.lock.Unlock()
  return self.dial()
}

// Dial will redial if conn is nil, and set
// timeouts if they've been set by the caller.
// 
// It simply returns nil if the socket appears already connected
func (self *ReusableConn)dial()(err os.Error){
  if self.conn == nil {
    self.conn, err = self.dialer()
    if err == nil && self.readTimeout != _UNSET_TIMEOUT {
      err = self.setReadTimeout(self.readTimeout)
    }
    if err == nil && self.writeTimeout != _UNSET_TIMEOUT {
      err = self.setWriteTimeout(self.writeTimeout)
    }
  }
  return
}

func (self *ReusableConn)close()(err os.Error){
  if self.conn != nil {
    err = self.conn.Close()
    self.conn = nil
  }
  return
}

func (self *ReusableConn)Close()(err os.Error){
  self.lock.Lock()
  defer self.lock.Unlock()
  return self.close()
}

// TODO: What's an appropriate responsde when we're not connected?
func (self *ReusableConn)RemoteAddr()(a net.Addr){
  self.lock.Lock()
  defer self.lock.Unlock()
  if self.conn != nil  {
    a = self.conn.RemoteAddr()
  }
  return
}

func (self *ReusableConn)LocalAddr()(a net.Addr){
  self.lock.Lock()
  defer self.lock.Unlock()
  if self.conn != nil  {
    a = self.conn.RemoteAddr()
  }
  return
}

func (self *ReusableConn)read(in []byte)(n int, err os.Error){
    err = self.dial()
    if err == nil {
      n, err = self.conn.Read(in)
      if err != nil { 
        self.close()
      }
    }
    return
}

func (self *ReusableConn)write(in []byte)(n int, err os.Error){
    err = self.dial()
    if err == nil {
      n, err = self.conn.Write(in)
      if err != nil { 
        self.close()
      }
    }
    return
}


func (self *ReusableConn)Read(in []byte)(n int, err os.Error){
  self.lock.Lock()
  defer self.lock.Unlock()
  return self.read(in)
}

func (self *ReusableConn)Write(out []byte)(n int, err os.Error){
  self.lock.Lock()
  defer self.lock.Unlock()
  return self.write(out)
}


func (self *ReusableConn)setReadTimeout(t int64)(err os.Error){
  err = self.dial()
  if err == nil {
    err = self.conn.SetReadTimeout(t)
    if err == nil {
      self.readTimeout = t
    }
  }
  return
}

func (self *ReusableConn)setWriteTimeout(t int64)(err os.Error){
  err = self.dial()
  if err == nil {
    err = self.conn.SetWriteTimeout(t)
    if err == nil {
      self.writeTimeout = t
    }
  }
  return
}

func (self *ReusableConn)SetReadTimeout(t int64)(err os.Error){
  self.lock.Lock()
  defer self.lock.Unlock()
  return self.setReadTimeout(t)
}

func (self *ReusableConn)SetWriteTimeout(t int64)(err os.Error){
  self.lock.Lock()
  defer self.lock.Unlock()
  return self.setWriteTimeout(t)
}

func (self *ReusableConn)SetTimeout(t int64)(err os.Error){
  err = self.SetReadTimeout(t)
  if err == nil {
    err = self.SetWriteTimeout(t)
  }
  return
}
