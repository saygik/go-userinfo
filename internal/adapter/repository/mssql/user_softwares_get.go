package mssql

func (r *Repository) GetUserSoftwares(user string) ([]int64, error) {
	var softwares []int64
	_, err := r.db.Select(&softwares, "GetUserSoftwares $1", user)
	return softwares, err
}

func (r *Repository) GetUserSoftwaresByFio(user string) ([]int64, error) {
	var softwares []int64
	_, err := r.db.Select(&softwares, "GetUserSoftwaresByFio $1", user)
	return softwares, err
}
