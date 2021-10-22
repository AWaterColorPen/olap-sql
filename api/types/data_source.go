package types

import (
	"fmt"
)

type DataSource struct {
	Name string `json:"name"`
}

func (d *DataSource) Statement() (string, error) {
	return fmt.Sprintf("%v", d.Name), nil
	// if d.SubRequest == nil {
	// 	return fmt.Sprintf("%v", d.Name), nil
	// }

	// statement, err := d.SubRequest.Statement()
	// if err != nil {
	// 	return "", err
	// }
	// return fmt.Sprintf("( %v ) %v", statement, d.Name), nil
}
