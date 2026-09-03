package web

import "github.com/gorilla/sessions"

type Session struct {
	c    *Context
	sess *sessions.Session
}

func (s *Session) Get(name string) string {
	if value, ok := s.sess.Values[name]; ok {
		return value.(string)
	}
	return ""
}

func (s *Session) GetDel(name string) string {
	if value, ok := s.sess.Values[name]; ok {
		s.Del(name)
		return value.(string)
	}
	return ""
}

func (s *Session) Exist(name string) bool {
	_, ok := s.sess.Values[name]
	return ok
}

func (s *Session) Set(name string, value any) {
	if value == nil {
		if !s.Exist(name) {
			return
		}
		// delete
		delete(s.sess.Values, name)
	} else {
		// set
		s.sess.Values[name] = value
	}
	s.sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 365,
		HttpOnly: true,
	}

	if err := s.sess.Save(s.c.Request.Http(), s.c.Response.Http()); err != nil {
		panic(err)
	}
}

func (s *Session) Del(name string) {
	s.Set(name, nil)
}
