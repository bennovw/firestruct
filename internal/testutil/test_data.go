package testutil

import (
	"time"

	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/type/latlng"
)

// TestSimpleStruct is a struct used for testing. All the fields names match the Firestore field names
type TestSimpleStruct struct {
	TimeData      time.Time
	StringData    string
	UUIDData      uuid.UUID
	BoolData      bool
	IntData       int
	DoubleData    float64
	BytesData     []byte
	NilData       any
	GeoPointData  latlng.LatLng
	ReferenceData string
	NestedMapData map[string]any
}

// TestTaggedStruct is a struct used for testing. All the fields names do not match the Firestore field names and are tagged as a result
type TestTaggedStruct struct {
	Time      time.Time      `firestore:"timeData"`
	String    string         `firestore:"stringData"`
	UUID      uuid.UUID      `firestore:"uuidData"`
	Bool      bool           `firestore:"boolData"`
	Int       int64          `firestore:"intData"`
	Double    float64        `firestore:"doubleData"`
	Bytes     []byte         `firestore:"bytesData"`
	Nil       any            `firestore:"nilData"`
	GeoPoint  latlng.LatLng  `firestore:"geoPointData"`
	Ref       string         `firestore:"referenceData"`
	NestedMap map[string]any `firestore:"nestedMapData"`
}

var TestFirebaseCloudEvents = []map[string]any{
	// FirestoreCloudEvent with no updateMask
	{
		"oldValue":   map[string]any{},
		"value":      TestFirebaseDocs[0],
		"updateMask": map[string]any{},
	},
}

var TestFirebaseDocs = []map[string]any{
	// Firestore document with a complex map field
	{
		"name":       "projects/projectID/databases/(default)/documents/collection/document",
		"fields":     TestFirebaseDocFields[12],
		"createTime": "2025-04-14T01:02:03Z",
		"updateTime": "2025-04-14T01:02:03Z",
	},
}

var TestFirebaseDocFields = []map[string]any{
	{
		"timeData": map[string]any{
			"timestampValue": testTimestamp,
		},
	},
	{
		"stringData": map[string]any{
			"stringValue": "Hello World",
		},
	},
	{
		"uuidData": map[string]any{
			"stringValue": "1f117a40-8bdb-4e8a-8f24-1622fea695b1",
		},
	},
	{
		"boolData": map[string]any{
			"booleanValue": true,
		},
	},
	{
		"intData": map[string]any{
			"integerValue": 987654321,
		},
	},
	{
		"doubleData": map[string]any{
			"doubleValue": 987.123456,
		},
	},
	{
		"bytesData": map[string]any{
			"bytesValue": []byte("Hello World"),
		},
	},
	{
		"nilData": map[string]any{
			"nullValue": nil,
		},
	},
	{
		"referenceData": map[string]any{
			"referenceValue": "/reference/path",
		},
	},
	{
		"geoPointData": map[string]any{
			"geoPointValue": latlng.LatLng{Latitude: 51.205005708080876, Longitude: 3.225345050850536},
		},
	},
	{
		"mapValue": map[string]any{
			"fields": map[string]any{
				"timeData": map[string]any{
					"timestampValue": testTimestamp,
				},
				"stringData": map[string]any{
					"stringValue": "Hello World",
				},
				"uuidData": map[string]any{
					"stringValue": "1f117a40-8bdb-4e8a-8f24-1622fea695b1",
				},
				"boolData": map[string]any{
					"booleanValue": true,
				},
				"intData": map[string]any{
					"integerValue": 987654321,
				},
				"doubleData": map[string]any{
					"doubleValue": 987.123456,
				},
				"bytesData": map[string]any{
					"bytesValue": []byte("Hello World"),
				},
				"nilData": map[string]any{
					"nullValue": nil,
				},
				"referenceData": map[string]any{
					"referenceValue": "/reference/path",
				},
				"geoPointData": map[string]any{
					"geoPointValue": latlng.LatLng{Latitude: 51.205005708080876, Longitude: 3.225345050850536},
				},
			},
		},
	},
	{
		"ArrayData": map[string]any{
			"arrayValue": map[string]any{
				"values": []any{

					map[string]any{
						"timeData": map[string]any{
							"timestampValue": testTimestamp,
						},
					},
					map[string]any{
						"stringData": map[string]any{
							"stringValue": "Hello World",
						},
					},
					map[string]any{
						"uuidData": map[string]any{
							"stringValue": "1f117a40-8bdb-4e8a-8f24-1622fea695b1",
						},
					},
					map[string]any{
						"boolData": map[string]any{
							"booleanValue": true,
						},
					},
					map[string]any{
						"intData": map[string]any{
							"integerValue": 987654321,
						},
					},
					map[string]any{
						"doubleData": map[string]any{
							"doubleValue": 987.123456,
						},
					},
					map[string]any{
						"bytesData": map[string]any{
							"bytesValue": []byte("Hello World"),
						},
					},
					map[string]any{
						"nilData": map[string]any{
							"nullValue": nil,
						},
					},
					map[string]any{
						"referenceData": map[string]any{
							"referenceValue": "/reference/path",
						},
					},
					map[string]any{
						"geoPointData": map[string]any{
							"geoPointValue": latlng.LatLng{Latitude: 51.205005708080876, Longitude: 3.225345050850536},
						},
					},
				},
			},
		},
	},
	// Test document [12] is a complex document containing all possible data types
	// it includes a nested map that contains an array containing both a map and another nested array
	{
		"timeData": map[string]any{
			"timestampValue": testTimestamp,
		},
		"stringData": map[string]any{
			"stringValue": "Hello World",
		},
		"uuidData": map[string]any{
			"stringValue": "1f117a40-8bdb-4e8a-8f24-1622fea695b1",
		},
		"boolData": map[string]any{
			"booleanValue": true,
		},
		"intData": map[string]any{
			"integerValue": "987654321",
		},
		"doubleData": map[string]any{
			"doubleValue": 987.123456,
		},
		"bytesData": map[string]any{
			"bytesValue": []byte("Hello World"),
		},
		"nilData": map[string]any{
			"nullValue": nil,
		},
		"referenceData": map[string]any{
			"referenceValue": "/reference/path",
		},
		"geoPointData": map[string]any{
			"geoPointValue": latlng.LatLng{Latitude: 51.205005708080876, Longitude: 3.225345050850536},
		},
		"nestedMapData": map[string]any{
			"mapValue": map[string]any{
				"fields": map[string]any{
					"nestedArrayData": map[string]any{
						"arrayValue": map[string]any{
							"values": []any{
								map[string]any{
									"mapValue": map[string]any{
										"fields": map[string]any{
											"timeData": map[string]any{
												"timestampValue": testTimestamp,
											},
											"stringData": map[string]any{
												"stringValue": "Hello World",
											},
											"uuidData": map[string]any{
												"stringValue": "1f117a40-8bdb-4e8a-8f24-1622fea695b1",
											},
											"boolData": map[string]any{
												"booleanValue": true,
											},
											"intData": map[string]any{
												"integerValue": 987654321,
											},
											"doubleData": map[string]any{
												"doubleValue": 987.123456,
											},
											"bytesData": map[string]any{
												"bytesValue": []byte("Hello World"),
											},
											"nilData": map[string]any{
												"nullValue": nil,
											},
											"referenceData": map[string]any{
												"referenceValue": "/reference/path",
											},
											"geoPointData": map[string]any{
												"geoPointValue": latlng.LatLng{Latitude: 51.205005708080876, Longitude: 3.225345050850536},
											},
										},
									},
								},
								map[string]any{
									"subNestedArrayData": map[string]any{
										"arrayValue": map[string]any{
											"values": []any{
												map[string]any{
													"timeData": map[string]any{
														"timestampValue": testTimestamp,
													},
												},
												map[string]any{
													"stringData": map[string]any{
														"stringValue": "Hello World",
													},
												},
												map[string]any{
													"uuidData": map[string]any{
														"stringValue": "1f117a40-8bdb-4e8a-8f24-1622fea695b1",
													},
												},
												map[string]any{
													"boolData": map[string]any{
														"booleanValue": true,
													},
												},
												map[string]any{
													"intData": map[string]any{
														"integerValue": 987654321,
													},
												},
												map[string]any{
													"doubleData": map[string]any{
														"doubleValue": 987.123456,
													},
												},
												map[string]any{
													"bytesData": map[string]any{
														"bytesValue": []byte("Hello World"),
													},
												},
												map[string]any{
													"nilData": map[string]any{
														"nullValue": nil,
													},
												},
												map[string]any{
													"referenceData": map[string]any{
														"referenceValue": "/reference/path",
													},
												},
												map[string]any{
													"geoPointData": map[string]any{
														"geoPointValue": latlng.LatLng{Latitude: 51.205005708080876, Longitude: 3.225345050850536},
													},
												},
											},
										},
									},
								},
								map[string]any{
									"timeData": map[string]any{
										"timestampValue": testTimestamp,
									},
								},
								map[string]any{
									"stringData": map[string]any{
										"stringValue": "Hello World",
									},
								},
								map[string]any{
									"uuidData": map[string]any{
										"stringValue": "1f117a40-8bdb-4e8a-8f24-1622fea695b1",
									},
								},
								map[string]any{
									"boolData": map[string]any{
										"booleanValue": true,
									},
								},
								map[string]any{
									"intData": map[string]any{
										"integerValue": 987654321,
									},
								},
								map[string]any{
									"doubleData": map[string]any{
										"doubleValue": 987.123456,
									},
								},
								map[string]any{
									"bytesData": map[string]any{
										"bytesValue": []byte("Hello World"),
									},
								},
								map[string]any{
									"nilData": map[string]any{
										"nullValue": nil,
									},
								},
								map[string]any{
									"referenceData": map[string]any{
										"referenceValue": "/reference/path",
									},
								},
								map[string]any{
									"geoPointData": map[string]any{
										"geoPointValue": latlng.LatLng{Latitude: 51.205005708080876, Longitude: 3.225345050850536},
									},
								},
							},
						},
					},
				},
			},
		},
	},
}

var FlattenedMapResults = []map[string]any{
	{
		"timeData": testTimestamp,
	},
	{
		"stringData": "Hello World",
	},
	{
		"uuidData": "1f117a40-8bdb-4e8a-8f24-1622fea695b1",
	},
	{
		"boolData": true,
	},
	{
		"intData": 987654321,
	},
	{
		"doubleData": 987.123456,
	},
	{
		"bytesData": []byte("Hello World"),
	},
	{
		"nilData": nil,
	},
	{
		"referenceData": "/reference/path",
	},
	{
		"geoPointData": latlng.LatLng{Latitude: 51.205005708080876, Longitude: 3.225345050850536},
	},
	{
		"timeData":      testTimestamp,
		"stringData":    "Hello World",
		"uuidData":      "1f117a40-8bdb-4e8a-8f24-1622fea695b1",
		"boolData":      true,
		"intData":       987654321,
		"doubleData":    987.123456,
		"bytesData":     []byte("Hello World"),
		"nilData":       nil,
		"referenceData": "/reference/path",
		"geoPointData":  latlng.LatLng{Latitude: 51.205005708080876, Longitude: 3.225345050850536},
	},
	{
		"ArrayData": []any{
			map[string]any{
				"timeData": testTimestamp,
			},
			map[string]any{
				"stringData": "Hello World",
			},
			map[string]any{
				"uuidData": "1f117a40-8bdb-4e8a-8f24-1622fea695b1",
			},
			map[string]any{
				"boolData": true,
			},
			map[string]any{
				"intData": 987654321,
			},
			map[string]any{
				"doubleData": 987.123456,
			},
			map[string]any{
				"bytesData": []byte("Hello World"),
			},
			map[string]any{
				"nilData": nil,
			},
			map[string]any{
				"referenceData": "/reference/path",
			},
			map[string]any{
				"geoPointData": latlng.LatLng{Latitude: 51.205005708080876, Longitude: 3.225345050850536},
			},
		},
	},

	{
		"timeData":      testTimestamp,
		"stringData":    "Hello World",
		"uuidData":      "1f117a40-8bdb-4e8a-8f24-1622fea695b1",
		"boolData":      true,
		"intData":       987654321,
		"doubleData":    987.123456,
		"bytesData":     []byte("Hello World"),
		"nilData":       nil,
		"referenceData": "/reference/path",
		"geoPointData":  latlng.LatLng{Latitude: 51.205005708080876, Longitude: 3.225345050850536},
		"nestedMapData": map[string]any{
			"nestedArrayData": []any{
				map[string]any{
					"timeData":      testTimestamp,
					"stringData":    "Hello World",
					"uuidData":      "1f117a40-8bdb-4e8a-8f24-1622fea695b1",
					"boolData":      true,
					"intData":       987654321,
					"doubleData":    987.123456,
					"bytesData":     []byte("Hello World"),
					"nilData":       nil,
					"referenceData": "/reference/path",
					"geoPointData":  latlng.LatLng{Latitude: 51.205005708080876, Longitude: 3.225345050850536},
				},
				map[string]any{
					"subNestedArrayData": []any{
						map[string]any{"timeData": testTimestamp},
						map[string]any{"stringData": "Hello World"},
						map[string]any{"uuidData": "1f117a40-8bdb-4e8a-8f24-1622fea695b1"},
						map[string]any{"boolData": true},
						map[string]any{"intData": 987654321},
						map[string]any{"doubleData": 987.123456},
						map[string]any{"bytesData": []byte("Hello World")},
						map[string]any{"nilData": nil},
						map[string]any{"referenceData": "/reference/path"},
						map[string]any{"geoPointData": latlng.LatLng{Latitude: 51.205005708080876, Longitude: 3.225345050850536}},
					},
				},
				map[string]any{
					"timeData": testTimestamp,
				},
				map[string]any{
					"stringData": "Hello World",
				},
				map[string]any{
					"uuidData": "1f117a40-8bdb-4e8a-8f24-1622fea695b1",
				},
				map[string]any{
					"boolData": true,
				},
				map[string]any{
					"intData": 987654321,
				},
				map[string]any{
					"doubleData": 987.123456,
				},
				map[string]any{
					"bytesData": []byte("Hello World"),
				},
				map[string]any{
					"nilData": nil,
				},
				map[string]any{
					"referenceData": "/reference/path",
				},
				map[string]any{
					"geoPointData": latlng.LatLng{Latitude: 51.205005708080876, Longitude: 3.225345050850536},
				},
			},
		},
	},
}

var testTimestamp, _ = time.Parse(time.RFC3339, "2025-04-14T01:02:03Z")
var testUUID, _ = uuid.Parse("1f117a40-8bdb-4e8a-8f24-1622fea695b1")
var testLatLang = latlng.LatLng{Latitude: 51.205005708080876, Longitude: 3.225345050850536}

var StructResults = []any{
	TestSimpleStruct{
		TimeData:      testTimestamp,
		StringData:    "Hello World",
		UUIDData:      testUUID,
		BoolData:      true,
		IntData:       987654321,
		DoubleData:    987.123456,
		BytesData:     []byte("Hello World"),
		NilData:       nil,
		GeoPointData:  testLatLang,
		ReferenceData: "/reference/path",
		NestedMapData: map[string]any{
			"nestedArrayData": []any{
				map[string]any{

					"timeData":      testTimestamp,
					"stringData":    "Hello World",
					"uuidData":      "1f117a40-8bdb-4e8a-8f24-1622fea695b1",
					"boolData":      true,
					"intData":       987654321,
					"doubleData":    987.123456,
					"bytesData":     []byte("Hello World"),
					"nilData":       nil,
					"referenceData": "/reference/path",
					"geoPointData":  latlng.LatLng{Latitude: 51.205005708080876, Longitude: 3.225345050850536},
				},
				map[string]any{
					"subNestedArrayData": []any{
						map[string]any{"timeData": testTimestamp},
						map[string]any{"stringData": "Hello World"},
						map[string]any{"uuidData": "1f117a40-8bdb-4e8a-8f24-1622fea695b1"},
						map[string]any{"boolData": true},
						map[string]any{"intData": 987654321},
						map[string]any{"doubleData": 987.123456},
						map[string]any{"bytesData": []byte("Hello World")},
						map[string]any{"nilData": nil},
						map[string]any{"referenceData": "/reference/path"},
						map[string]any{"geoPointData": latlng.LatLng{Latitude: 51.205005708080876, Longitude: 3.225345050850536}},
					},
				},
				map[string]any{
					"timeData": testTimestamp,
				},
				map[string]any{
					"stringData": "Hello World",
				},
				map[string]any{
					"uuidData": "1f117a40-8bdb-4e8a-8f24-1622fea695b1",
				},
				map[string]any{
					"boolData": true,
				},
				map[string]any{
					"intData": 987654321,
				},
				map[string]any{
					"doubleData": 987.123456,
				},
				map[string]any{
					"bytesData": []byte("Hello World"),
				},
				map[string]any{
					"nilData": nil,
				},
				map[string]any{
					"referenceData": "/reference/path",
				},
				map[string]any{
					"geoPointData": latlng.LatLng{Latitude: 51.205005708080876, Longitude: 3.225345050850536},
				},
			},
		},
	},
	TestTaggedStruct{
		Time:     testTimestamp,
		String:   "Hello World",
		UUID:     testUUID,
		Bool:     true,
		Int:      987654321,
		Double:   987.123456,
		Bytes:    []byte("Hello World"),
		Nil:      nil,
		GeoPoint: testLatLang,
		Ref:      "/reference/path",
		NestedMap: map[string]any{
			"nestedArrayData": []any{
				map[string]any{

					"timeData":      testTimestamp,
					"stringData":    "Hello World",
					"uuidData":      "1f117a40-8bdb-4e8a-8f24-1622fea695b1",
					"boolData":      true,
					"intData":       987654321,
					"doubleData":    987.123456,
					"bytesData":     []byte("Hello World"),
					"nilData":       nil,
					"referenceData": "/reference/path",
					"geoPointData":  latlng.LatLng{Latitude: 51.205005708080876, Longitude: 3.225345050850536},
				},
				map[string]any{
					"subNestedArrayData": []any{
						map[string]any{"timeData": testTimestamp},
						map[string]any{"stringData": "Hello World"},
						map[string]any{"uuidData": "1f117a40-8bdb-4e8a-8f24-1622fea695b1"},
						map[string]any{"boolData": true},
						map[string]any{"intData": 987654321},
						map[string]any{"doubleData": 987.123456},
						map[string]any{"bytesData": []byte("Hello World")},
						map[string]any{"nilData": nil},
						map[string]any{"referenceData": "/reference/path"},
						map[string]any{"geoPointData": latlng.LatLng{Latitude: 51.205005708080876, Longitude: 3.225345050850536}},
					},
				},
				map[string]any{
					"timeData": testTimestamp,
				},
				map[string]any{
					"stringData": "Hello World",
				},
				map[string]any{
					"uuidData": "1f117a40-8bdb-4e8a-8f24-1622fea695b1",
				},
				map[string]any{
					"boolData": true,
				},
				map[string]any{
					"intData": 987654321,
				},
				map[string]any{
					"doubleData": 987.123456,
				},
				map[string]any{
					"bytesData": []byte("Hello World"),
				},
				map[string]any{
					"nilData": nil,
				},
				map[string]any{
					"referenceData": "/reference/path",
				},
				map[string]any{
					"geoPointData": latlng.LatLng{Latitude: 51.205005708080876, Longitude: 3.225345050850536},
				},
			},
		},
	},
}
