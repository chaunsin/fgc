package parse

import (
	"fmt"
	"os"
	"path/filepath"
)

func readOrder(dir string, cc *CryptoConfig) error {
	dir = filepath.Join(dir, orderPath)
	list, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("ReadDir:%w", err)
	}

	for _, v := range list {
		if !v.IsDir() {
			continue
		}

		var (
			org Org
			p   = filepath.Join(dir, v.Name())
		)
		cc.Order[OrgName(v.Name())] = &org

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
				if err = msp(p, oo, &m); err == nil {
					cc.Order[OrgName(v.Name())].Msp = m
				}
			case "orderers":
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
		fmt.Printf("readOrder:%+v\n", org)
	}

	return nil
}
