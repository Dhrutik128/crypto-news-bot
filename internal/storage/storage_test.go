package storage

import (
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/buntdb"
	tb "gopkg.in/tucnak/telebot.v2"
	"testing"
)

func TestDB_Set(t *testing.T) {
	type fields struct {
		DB *buntdb.DB
	}
	type args struct {
		object User
	}
	dbtest, err := buntdb.Open("../../data/data_test.db")
	if err != nil {
		log.Fatal(err)
	}
	err = dbtest.CreateIndex("user", "user_*", buntdb.IndexJSON("user.id"))
	if err != nil {
		panic(err)
	}
	log.Infoln("started database")
	//database := &DB{DB: db}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "t1", args: args{object: User{User: &tb.User{ID: 12312312, FirstName: "tester", LastName: "one"}}}, fields: fields{DB: dbtest}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &DB{
				DB: tt.fields.DB,
			}
			if err := db.Set(tt.args.object); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			user := User{User: &tb.User{ID: 12312312}}
			if err := db.Get(&user); err != nil {
				t.Errorf("Get() error = %v", err)
				return
			}
			if user.User.FirstName != tt.args.object.User.FirstName {
				t.Errorf("Firstname error")
				return
			}
			u := User{User: &tb.User{ID: 12312312}}
			if err := db.Delete(u); err != nil {
				t.Errorf("Delete() error = %v", err)
				return
			}
			if err := db.Get(&u); err == nil {
				t.Errorf("Get() found deleted entry")
				return
			}
		})
	}
}
