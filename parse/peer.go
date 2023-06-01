package parse

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ca(parent string, dir os.DirEntry, org *Org) error {
	var (
		pkg = Package{
			parentDir:  parent,
			currentDir: "ca",
		}
		p = filepath.Join(parent, "ca")
	)

	ca, err := os.ReadDir(p)
	if err != nil {
		return fmt.Errorf("ReadDir:%w", err)
	}
	if len(ca) < 2 {
		return fmt.Errorf("ca dir file lt < 0")
	}
	for _, c := range ca {
		if c.IsDir() {
			continue
		}
		path := File(filepath.Join(p, c.Name()))
		if strings.HasSuffix(c.Name(), ".pem") {
			pkg.CA = path
			continue
		}
		if strings.HasSuffix(c.Name(), "_sk") {
			pkg.Key = path
			continue
		}
		fmt.Printf("ca/%s is un used\n", c.Name())
	}
	org.CA = pkg
	return nil
}

func msp(parent string, dir os.DirEntry, m *Msp) error {
	var mspDir = filepath.Join(parent, "msp")
	msp, err := os.ReadDir(mspDir)
	if err != nil {
		return fmt.Errorf("ReadDir:%w", err)
	}
	if len(msp) < 3 {
		return fmt.Errorf("msp dir file lt < 3")
	}
	for _, v := range msp {
		path := filepath.Join(mspDir, v.Name())
		if !v.IsDir() {
			if v.Name() == "config.yaml" {
				m.ConfigYaml = Package{
					parentDir: "msp",
					Yaml:      File(path),
				}
			}
			continue
		}
		switch v.Name() {
		case "admincerts":
			{
				pkg := Package{
					parentDir:  "msp",
					currentDir: "admincerts",
				}
				d, err := os.ReadDir(filepath.Join(mspDir, v.Name()))
				if err != nil {
					return fmt.Errorf("ReadDir:%w", err)
				}
				if len(d) <= 0 {
					//fmt.Println("msp/admincerts file <=0")
					continue
				}
				for _, v := range d {
					if v.IsDir() {
						return fmt.Errorf("msp/admincerts %s is not file", v.Name())
					}
					if !strings.HasSuffix(v.Name(), ".pem") {
						return fmt.Errorf("msp/admincerts %s is valid file name", v.Name())
					}
					pkg.Cert = File(filepath.Join(path, v.Name()))
				}
				m.AdminCerts = pkg
			}
		case "cacerts":
			{
				pkg := Package{
					parentDir:  "msp",
					currentDir: "cacerts",
				}
				d, err := os.ReadDir(filepath.Join(mspDir, v.Name()))
				if err != nil {
					return fmt.Errorf("ReadDir:%w", err)
				}
				if len(d) <= 0 {
					return fmt.Errorf("msp/cacerts file <= 0")
				}
				for _, v := range d {
					if v.IsDir() {
						return fmt.Errorf("msp/cacerts %s is not file", v.Name())
					}
					if !strings.HasSuffix(v.Name(), ".pem") {
						return fmt.Errorf("msp/cacerts %s is valid file name", v.Name())
					}
					pkg.Cert = File(filepath.Join(path, v.Name()))
				}
				m.CaCerts = pkg
			}
		case "tlscacerts":
			{
				pkg := Package{
					parentDir:  "msp",
					currentDir: "tlscacerts",
				}
				d, err := os.ReadDir(filepath.Join(mspDir, v.Name()))
				if err != nil {
					return fmt.Errorf("ReadDir:%w", err)
				}
				if len(d) <= 0 {
					return fmt.Errorf("msp/tlscacerts file <= 0")
				}
				for _, v := range d {
					if v.IsDir() {
						return fmt.Errorf("msp/tlscacerts %s is not file", v.Name())
					}
					if !strings.HasSuffix(v.Name(), ".pem") {
						return fmt.Errorf("msp/tlscacerts %s is valid file name", v.Name())
					}
					pkg.Cert = File(filepath.Join(path, v.Name()))
				}
				m.TLSCaCerts = pkg
			}
		case "keystore":
			{
				pkg := Package{
					parentDir:  "msp",
					currentDir: "keystore",
				}
				d, err := os.ReadDir(filepath.Join(mspDir, v.Name()))
				if err != nil {
					return fmt.Errorf("ReadDir:%w", err)
				}
				if len(d) <= 0 {
					return fmt.Errorf("msp/keystore file lt <= 0")
				}
				for _, v := range d {
					if v.IsDir() {
						return fmt.Errorf("msp/keystore %s is not file", v.Name())
					}
					if !strings.HasSuffix(v.Name(), "_sk") {
						return fmt.Errorf("msp/keystore %s is valid file name", v.Name())
					}
					pkg.Key = File(filepath.Join(path, v.Name()))
				}
				m.KeyStore = pkg
			}
		case "signcerts":
			{
				pkg := Package{
					parentDir:  "msp",
					currentDir: "signcerts",
				}
				d, err := os.ReadDir(filepath.Join(mspDir, v.Name()))
				if err != nil {
					return fmt.Errorf("ReadDir:%w", err)
				}
				if len(d) <= 0 {
					return fmt.Errorf("msp/sigcerts file <= 0")
				}
				for _, v := range d {
					if v.IsDir() {
						return fmt.Errorf("msp/sigcerts %s is not file", v.Name())
					}
					if !strings.HasSuffix(v.Name(), ".pem") {
						return fmt.Errorf("msp/sigcerts %s is valid file name", v.Name())
					}
					pkg.Cert = File(filepath.Join(path, v.Name()))
				}
				m.SignCerts = pkg
			}
		default:
			fmt.Printf("msp/%s is not used\n", v.Name())
		}
	}
	//fmt.Printf("msp:%+v\n", m)
	return nil
}

// server 读取peer或order
func server(parent string, dir os.DirEntry, org *Org) error {
	org.Server = map[OrgDomain]*Serve{}
	var opDir = filepath.Join(parent, dir.Name())
	m, err := os.ReadDir(opDir)
	if err != nil {
		return fmt.Errorf("ReadDir:%w", err)
	}
	if len(m) <= 0 {
		return fmt.Errorf("%s dir file lte < 0", dir.Name())
	}
	for _, v := range m {
		if !v.IsDir() {
			continue
		}
		orgName := v.Name()
		if _, ok := org.Server[OrgDomain(orgName)]; !ok {
			org.Server[OrgDomain(orgName)] = &Serve{}
		}
		p := filepath.Join(opDir, v.Name())
		d, err := os.ReadDir(p)
		if err != nil {
			return fmt.Errorf("ReadDir:%w", err)
		}
		if len(d) < 2 {
			return fmt.Errorf("%s file lt < 0", p)
		}
		for _, v := range d {
			if !v.IsDir() {
				continue
			}
			switch v.Name() {
			case "msp":
				var m Msp
				if err := msp(p, v, &m); err != nil {
					return fmt.Errorf("msp:%w", err)
				}
				org.Server[OrgDomain(orgName)].Msp = m
			case "tls":
				{
					path := filepath.Join(p, v.Name())
					d, err := os.ReadDir(path)
					if err != nil {
						return fmt.Errorf("ReadDir:%w", err)
					}
					if len(d) < 3 {
						return fmt.Errorf("%s/tls file lt < 3", d)
					}
					for _, v := range d {
						if v.IsDir() {
							continue
						}
						switch v.Name() {
						case "ca.crt":
							org.Server[OrgDomain(orgName)].TLS.CA = File(filepath.Join(path, v.Name()))
						case "server.crt":
							org.Server[OrgDomain(orgName)].TLS.Cert = File(filepath.Join(path, v.Name()))
						case "server.key":
							org.Server[OrgDomain(orgName)].TLS.Key = File(filepath.Join(path, v.Name()))
						default:
							fmt.Printf("%s un used\n", v.Name())
						}
					}
				}
			}
		}
	}
	return nil
}

func tlsCa(parent string, dir os.DirEntry, org *Org) error {
	var (
		tlsDir = filepath.Join(parent, "tlsca")
		pkg    = Package{
			parentDir:  parent,
			currentDir: "tls",
		}
	)
	m, err := os.ReadDir(tlsDir)
	if err != nil {
		return err
	}
	if len(m) < 2 {
		return fmt.Errorf("tlsca dir file lte < 2")
	}
	for _, v := range m {
		if v.IsDir() {
			continue
		}
		path := File(filepath.Join(tlsDir, v.Name()))
		if strings.HasSuffix(v.Name(), "_sk") {
			pkg.Key = path
			continue
		}
		if strings.HasSuffix(v.Name(), ".pem") {
			pkg.Cert = path
			continue
		}
		log.Printf("tlsca/%s is un used\n", v.Name())
	}
	org.TLSCA = pkg
	return nil
}

func users(parent string, dir os.DirEntry, org *Org) error {
	org.Users = map[UserDomain]*User{}
	var mspDir = filepath.Join(parent, "users")
	m, err := os.ReadDir(mspDir)
	if err != nil {
		return fmt.Errorf("ReadDir:%w", err)
	}
	if len(m) <= 0 {
		return fmt.Errorf("users dir file lte < 0")
	}
	for _, v := range m {
		if !v.IsDir() {
			continue
		}
		username := v.Name()
		if _, ok := org.Users[UserDomain(username)]; !ok {
			org.Users[UserDomain(username)] = &User{}
		}
		p := filepath.Join(mspDir, v.Name())
		d, err := os.ReadDir(p)
		if err != nil {
			return fmt.Errorf("ReadDir:%w", err)
		}
		if len(d) < 2 {
			return fmt.Errorf("msp/tlscacerts file lt < 0")
		}
		for _, v := range d {
			if !v.IsDir() {
				continue
			}
			switch v.Name() {
			case "msp":
				var m Msp
				if err := msp(p, v, &m); err != nil {
					return fmt.Errorf("msp:%w", err)
				}
				org.Users[UserDomain(username)].Msp = m
			case "tls":
				{
					path := filepath.Join(p, v.Name())
					d, err := os.ReadDir(path)
					if err != nil {
						return fmt.Errorf("ReadDir:%w", err)
					}
					if len(d) < 3 {
						return fmt.Errorf("%s/tls file lt < 3", v)
					}
					for _, v := range d {
						if v.IsDir() {
							continue
						}
						cur := File(filepath.Join(path, v.Name()))
						switch v.Name() {
						case "ca.crt":
							org.Users[UserDomain(username)].TLS.CA = cur
						case "client.crt":
							org.Users[UserDomain(username)].TLS.Cert = cur
						case "client.key":
							org.Users[UserDomain(username)].TLS.Key = cur
						default:
							fmt.Printf("%s un used\n", v.Name())
						}
					}
				}
			}
		}
	}
	return nil
}

func readPeer(dir string, cc *CryptoConfig) error {
	dir = filepath.Join(dir, peerPath)
	list, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("ReadDir():%w", err)
	}

	for _, v := range list {
		if !v.IsDir() {
			continue
		}

		var (
			org Org
			p   = filepath.Join(dir, v.Name())
		)
		cc.Orgs[OrgName(v.Name())] = &org

		o, err := os.ReadDir(p)
		if err != nil {
			return fmt.Errorf("ReadDir:%w", err)
		}
		if len(o) < 5 {
			return fmt.Errorf("peer org is < 5")
		}
		for _, oo := range o {
			switch oo.Name() {
			case "ca":
				err = ca(p, oo, &org)
			case "msp":
				var m Msp
				if err = msp(p, v, &m); err == nil {
					cc.Orgs[OrgName(v.Name())].Msp = m
				}
			case "peers":
				err = server(p, oo, &org)
			case "tlsca":
				err = tlsCa(p, oo, &org)
			case "users":
				err = users(p, oo, &org)
			default:
				fmt.Printf("1. %s => %s unuse file\n", p, v.Name())
			}
			if err != nil {
				return fmt.Errorf("org:%w", err)
			}
		}
		fmt.Printf("readPeer:%+v\n", org)
	}
	return nil
}
