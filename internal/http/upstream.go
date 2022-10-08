/*
 *  Copyright (c) 2022 NetEase Inc.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

/*
 * Project: Pigeon
 * Created Date: 2022-09-20
 * Author: Jingli Chen (Wine93)
 */

package http


/*
import (
	"net"
	"sync"
	"time"
)

const (
	STATUS_OK = iota
	STATUS_ERROR
)

type Server struct {
	sync.RWMutex    // lock status
	Address         string
	Weight          int
	EffectiveWeight int
	CurrentWeight   int
	MaxFails        int `default:"3"`
	FailTimeout     int
	LastModify      time.Duration
	Status          int `default:"0"`
}

type Config struct {
	CheckInterval time.Duration
	CheckTimeout  time.Duration
}

type Upstream struct {
	servers []Server
	cfg     *Config
}

func NewUpstream() *Upstream {
	return &Upstream{}
}

func (u *Upstream) connect(s *Server, wg *sync.WaitGroup) {
	defer wg.Done()
	_, err := net.DialTimeout("tcp", s.Address, time.Duration(3)*time.Second)
	if err == nil {
		return
	}

	s.Lock()
	defer s.Unlock()
	s.Status = STATUS_ERROR
}

func (u *Upstream) checkHealth() {
	var wg sync.WaitGroup
	for _, server := range u.servers {
		wg.Add(1)
		go u.connect(&server, &wg)
	}
	wg.Wait()
}

func (u *Upstream) StartChecker() {
	timer := time.NewTimer(u.cfg.CheckInterval)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			u.checkHealth()
			timer.Reset(u.cfg.CheckInterval)
		}
	}
}

func (u *Upstream) isHealthy(s *Server) bool {
	s.RLock()
	defer s.RUnlock()
	return s.Status == STATUS_OK
}

func (u *Upstream) nextRoundRobinServer() (*Server, error) {
	return nil, nil
}

func (u *Upstream) Get() (*Server, error) {
	return u.nextRoundRobinServer()
}
 */
