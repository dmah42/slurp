package main

import (
	"log"

	"github.com/alexanderritola/nntp"
)

type NNTP struct {
	// map from address to connection
	conns map[string]*nntp.Conn
}

func NewNNTP(c *Config) (*NNTP, error) {
	n := &NNTP{
		conns: make(map[string]*nntp.Conn),
	}
	for nw := range c.Networks {
		network := c.Networks[nw]
		addr := network.Address

		// dial
		// TODO: TLS
		log.Printf("connecting to %q\n", addr)
		conn, err := nntp.Dial("tcp", addr)
		if err != nil {
			return nil, err
		}

		// auth
		log.Printf("authenticating as %q\n", network.User)
		if err := conn.Authenticate(network.User, network.Pass); err != nil {
			conn.Quit()
			return nil, err
		}
		n.conns[addr] = conn
	}
	return n, nil
}

func (n *NNTP) Addresses() ([]string, error) {
	as := make([]string, len(n.conns))
	i := 0
	for a := range n.conns {
		as[i] = a
		i += 1
	}
	return as, nil
}

func (n *NNTP) Groups(addr string) ([]*nntp.Group, error) {
	log.Printf("listing groups on %q\n", addr)
	gs, err := n.conns[addr].List()
	if err != nil {
		return nil, err
	}
	return gs, nil
}

func (n *NNTP) Articles(addr string, group string) ([]nntp.MessageOverview, error) {
	log.Printf("switching to %q on %q\n", group, addr)
	g, err := n.conns[addr].Group(group)
	if err != nil {
		return nil, err
	}

	log.Printf("got %d articles", g.High-g.Low)
	max := g.High
	if g.High-g.Low > 100 {
		max = g.Low + 100
	}

	log.Printf("listing %d articles\n", max-g.Low)
	ovs, err := n.conns[addr].Overview(g.Low, max)
	if err != nil {
		return nil, err
	}

	return ovs, err
}

/* TODO
// connect to a news group
	grp := "alt.binaries.pictures"
	_, l, _, err := conn.Group(grp)
	if err != nil {
		log.Fatalf("Could not connect to group %s: %v %d", grp, err, l)
	}

	// fetch an article
	id := "<4c1c18ec$0$8490$c3e8da3@news.astraweb.com>"
	article, err := conn.Article(id)
	if err != nil {
		log.Fatalf("Could not fetch article %s: %v", id, err)
	}

	// read the article contents
	body, err := ioutil.ReadAll(article.Body)
	if err != nil {
		log.Fatalf("error reading reader: %v", err)
	}
*/

func (n *NNTP) Close() []error {
	var es []error
	for c := range n.conns {
		e := n.conns[c].Quit()
		if e != nil {
			es = append(es, e)
		}
	}
	return es
}
