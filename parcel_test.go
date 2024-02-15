package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	randSource = rand.NewSource(time.Now().UnixNano())
	randRange  = rand.New(randSource)
)

func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

func TestAddGetDelete(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")

	store := NewParcelStore(db)
	parcel := getTestParcel()
	id, err := store.Add(parcel) //adding new parcel in db
	require.NoError(t, err)      //check
	require.NotEmpty(t, id)
	parcel.Number = id      //number assignment after adding the database entry
	p, err := store.Get(id) //check
	assert.Equal(t, parcel, p)
	require.NoError(t, err)
	err = store.Delete(id) //deleting this db entry
	require.NoError(t, err)
	p, err = store.Get(id) //check
	assert.Error(t, err)
	assert.Empty(t, p)
}

func TestSetAddress(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")

	require.NoError(t, err)
	store := NewParcelStore(db)
	parcel := getTestParcel()
	id, err := store.Add(parcel) //adding new parcel in db
	require.NoError(t, err)
	require.NotEmpty(t, id)
	newAddress := "new test address"
	err = store.SetAddress(id, newAddress) //setting new address
	require.NoError(t, err)
	p, err := store.Get(id) //check
	require.Equal(t, newAddress, p.Address)
	require.NoError(t, err)
}

func TestSetStatus(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")

	require.NoError(t, err)
	store := NewParcelStore(db)
	parcel := getTestParcel()
	id, err := store.Add(parcel) //adding new parcel in db
	require.NoError(t, err)
	require.Greater(t, id, 0)
	require.NotEmpty(t, id)
	newStatus := ParcelStatusSent
	err = store.SetStatus(id, newStatus) //setting new status and checking for errors
	require.NoError(t, err)
	p, err := store.Get(id)
	require.Equal(t, newStatus, p.Status)
}

func TestGetByClient(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")

	require.NoError(t, err)
	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelSl := []Parcel{} //replacing map with slice type to use assert.Containsf()

	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	// add
	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		require.NotEmpty(t, id)

		parcels[i].Number = id
		parcelSl = append(parcelSl, parcels[i])
	}

	// get by client
	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	require.Equal(t, len(parcels), len(storedParcels))

	// check
	for _, parcel := range storedParcels {
		assert.Containsf(t, parcelSl, parcel, "Slice parcelSl does not conatain parcel")
	}
}
