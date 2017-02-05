package chipmunk

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB struct {
	db *sqlx.DB
}

func NewDB(dbhost, dbname string) (*DB, error) {
	var err error
	config := fmt.Sprintf(
		"dbname=%s host=%s sslmode=disable",
		dbname,
		dbhost,
	)
	db, err := sqlx.Connect(
		"postgres",
		config,
	)
	if err != nil {
		return nil, err
	}

	d := &DB{db}
	d.initializeDB()
	return d, nil
}

func (d *DB) initializeDB() {
	// XXX ignoring errors
	d.db.Exec(createdb)
	d.db.Exec(primeCategories)
}

func (d *DB) getCategories() ([]category, error) {
	results := []category{}
	rows, err := d.db.Queryx("SELECT id, name, budget FROM categories")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var result category
		err := rows.StructScan(&result)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

//func (d *DB) checkOwner(filename, client string) error {
//	row := d.db.QueryRowx("SELECT client FROM pics WHERE filename = $1", filename)
//	var owner string
//	err := row.Scan(&owner)
//	if err == sql.ErrNoRows {
//		return errors.New("file not in DB")
//	}
//	if owner != strings.ToLower(client) {
//		return errors.New("filename did not match owner")
//	}
//	return nil
//}
//
//func (d *DB) getPic(filename string) (Pic, error) {
//	p := Pic{}
//	err := d.db.QueryRowx(
//		`
//		SELECT
//			filename,
//			category_id as category,
//			client,
//			categorized,
//			reported
//		FROM pics
//		WHERE
//			filename = $1
//		LIMIT 1
//		`,
//		filename,
//	).StructScan(&p)
//	logs.Debug("%+v", p)
//	if err == sql.ErrNoRows {
//		return Pic{}, errors.New("pic not found")
//	}
//	return p, nil
//}
//
//func (d *DB) votesForPic(filename string) ([]int, error) {
//	results := []int{}
//	rows, err := d.db.Queryx(`SELECT vote FROM votes WHERE filename = $1`, filename)
//	if err != nil {
//		return []int{}, err
//	}
//	for rows.Next() {
//		var result int
//		err := rows.Scan(&result)
//		if err != nil {
//			return []int{}, err
//		}
//		results = append(results, result)
//	}
//	return results, nil
//}
//
//func (d *DB) closeVoting(filename string) error {
//	_, err := d.db.Exec(
//		`UPDATE pics SET categorized = true WHERE filename = $1`,
//		filename,
//	)
//	return err
//}
//
//func (d *DB) votingClosed(filename string) (bool, error) {
//	row := d.db.QueryRow("SELECT categorized FROM pics WHERE filename = $1", filename)
//	already := false
//	err := row.Scan(&already)
//	return already, err
//}
//
//type Token string
//
//func (d *DB) tokenForFile(filename string) (Token, error) {
//	row := d.db.QueryRow(`
//	SELECT
//		clients.token
//	FROM pics, clients
//	WHERE
//		pics.client = clients.id
//		AND pics.filename = $1
//	`,
//		filename,
//	)
//	token := sql.NullString{}
//	err := row.Scan(&token)
//	return Token(token.String), err
//}
//
//type Pic struct {
//	Filename    string `json:"filename"`
//	Category    int    `json:"category"`
//	Client      string `json:"client"`
//	Categorized bool   `json:"categorized"db:"categorized"`
//	Reported    bool   `json:"reported"db:"reported"`
//}
//
//type PicDetails struct {
//	// did the request work; this goes along with HTTP status codes
//	Success bool `json:"success"`
//
//	Filename string `json:"filename"`
//
//	Message     string `json:"message"`
//	Works       bool   `json:"works"`
//	Categorized bool   `json:"categorized"`
//}
//
//func (d *DB) getPicDetails(filename string) (PicDetails, error) {
//	pd := PicDetails{
//		Success:  false,
//		Filename: filename,
//		Message:  "unknown ranking",
//	}
//	votes, err := d.votesForPic(filename)
//	if err != nil {
//		return pd, err
//	}
//	message, works, done, err := categorize(votes)
//	if err != nil {
//		return pd, err
//	}
//	pd.Success = true
//	pd.Message = message
//	pd.Works = works
//	pd.Categorized = done
//	return pd, nil
//}
//
//type VoteReq struct {
//	Filename string `json:"filename"`
//	Client   string `json:"client"`
//	Value    int    `json:"value"`
//}
//
//type Category struct {
//	Id   int    `json:"id"db:"id"`
//	Name string `json:"name"db:"name"`
//}
//
//
//func (d *DB) flagPic(filename string, reporter string) (int, error) {
//	logs.Debug("(%s, %s)", filename, reporter)
//
//	var count int
//	row := d.db.QueryRow(
//		"SELECT COUNT(*) FROM reports WHERE filename = $1 AND reporter = $2",
//		filename,
//		reporter,
//	)
//	err := row.Scan(&count)
//	if err != nil {
//		return http.StatusInternalServerError, err
//	}
//	if count > 0 {
//		return http.StatusBadRequest, fmt.Errorf("%s already reported %s", reporter, filename)
//	}
//
//	_, err = d.db.Exec(
//		`INSERT INTO reports (filename, reporter) VALUES ($1, $2)`,
//		filename,
//		reporter,
//	)
//	if err != nil {
//		return http.StatusBadRequest, err
//	}
//
//	var culprit string
//	err = d.db.Get(&culprit, `
//		SELECT
//			client
//		FROM pics, clients
//		WHERE
//			pics.client = clients.id AND
//			pics.filename = $1
//		`,
//		filename,
//	)
//	if err != nil {
//		return http.StatusInternalServerError, err
//	}
//
//	_, err = d.db.Exec(`
//		UPDATE
//			clients
//		SET
//			infractions = infractions + 1
//		WHERE
//			clients.id = $1
//		`,
//		culprit,
//	)
//	if err != nil {
//		return http.StatusInternalServerError, err
//	}
//
//	_, err = d.db.Exec(`
//		UPDATE
//			pics
//		SET
//			reported = true
//		WHERE
//			filename = $1
//		`,
//		filename,
//	)
//	if err != nil {
//		return http.StatusInternalServerError, err
//	}
//	return http.StatusCreated, nil
//}
//
//type client struct {
//	Id          string
//	Token       sql.NullString
//	Infractions int
//	db          *sqlx.DB
//}
//
//const reportThreshold = 3
//
//func (c *client) isBlocked() bool {
//	if c.Infractions >= reportThreshold {
//		return true
//	}
//	return false
//}
//
//func (c *client) MarshalJSON() ([]byte, error) {
//	r := struct {
//		Id          string `json:"id"`
//		Infractions int    `json:"infractions"`
//	}{
//		Id:          c.Id,
//		Infractions: c.Infractions,
//	}
//
//	return json.Marshal(r)
//}
//
//func (c *client) SetToken(token string) error {
//	_, err := c.db.Exec(
//		`UPDATE clients SET token = $1 WHERE id = $2`,
//		token,
//		c.Id,
//	)
//	if err != nil {
//		return err
//	}
//	c.Token = sql.NullString{
//		String: token,
//		Valid:  true,
//	}
//	return nil
//}
//
//// getOrCreateClient performs db operations requisite to fetch current info or
//// create a new client.
////
//// returns:
//// *client, created or not, error
//func (d *DB) getOrCreateClient(clientId string) (*client, bool, error) {
//	var created bool
//	c := &client{
//		db: d.db,
//	}
//	err := d.db.QueryRowx(
//		"SELECT * FROM clients WHERE id = $1",
//		clientId,
//	).StructScan(c)
//	if err == sql.ErrNoRows {
//		_, err = d.db.Exec(
//			`INSERT INTO clients (id) VALUES ($1)`,
//			clientId,
//		)
//		if err != nil {
//			return nil, false, fmt.Errorf(
//				"problem inserting client: %s: %+v",
//				clientId,
//				err,
//			)
//		}
//		created = true
//		err := d.db.QueryRowx(
//			"SELECT * FROM clients WHERE id = $1",
//			clientId,
//		).StructScan(c)
//		if err != nil {
//			return nil, true, err
//		}
//	}
//	return c, created, nil
//}
//
//type UploadResp struct {
//	Success bool   `json:"success"`
//	Sha     string `json:"sha1"`
//}
