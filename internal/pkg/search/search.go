package search

import (
	"errors"
	"fmt"
	"game-server-selector/internal/models"
	"reflect"
	"strings"
)

type SearchService interface {
	//SearchAll(req []models.SearchRequest, servers []models.Server) ([]models.Server, error)
	SearchOne(req []models.SearchRequest, server models.Server) (bool, error)
}
type searchService struct {
	methods map[string]func(req models.SearchRequest, data models.Server) (bool, error)
}

const (
	OperatorEqual              = "eq"
	OperatorNotEqual           = "ne"
	OperatorGreaterThan        = "gt"
	OperatorGreaterThanOrEqual = "gte"
	OperatorLessThan           = "lt"
	OperatorLessThanOrEqual    = "lte"
	OperatorMatch              = "match"
	OperatorIn                 = "in"
)

func NewSearchService() SearchService {
	s := &searchService{
		methods: map[string]func(req models.SearchRequest, data models.Server) (bool, error){},
	}
	s.RegisterMethod(OperatorEqual, s.EqSearch)
	s.RegisterMethod(OperatorNotEqual, s.NeSearch)
	s.RegisterMethod(OperatorMatch, s.MatchSearch)
	s.RegisterMethod(OperatorIn, s.InSearch)
	s.RegisterMethod(OperatorGreaterThan, s.GtSearch)
	s.RegisterMethod(OperatorGreaterThanOrEqual, s.GteSearch)
	s.RegisterMethod(OperatorLessThan, s.LtSearch)
	s.RegisterMethod(OperatorLessThanOrEqual, s.LteSearch)
	return s
}

func (s *searchService) RegisterMethod(method string, f func(req models.SearchRequest, data models.Server) (bool, error)) {
	s.methods[method] = f
}

func (s *searchService) SearchOne(req []models.SearchRequest, server models.Server) (bool, error) {
	for _, v := range req {
		if f, ok := s.methods[v.Query.Operator]; ok {
			result, err := f(v, server)
			if err != nil {
				return false, err
			}
			if !result {
				return false, nil
			}
		}
	}
	return true, nil
}

func (s *searchService) EqSearch(req models.SearchRequest, data models.Server) (bool, error) {
	fieldVal, err := getFieldValue(data, req.Field)
	if err != nil {
		return false, err
	}
	return reflect.DeepEqual(fieldVal, req.Query.Value), nil
}

func (s *searchService) NeSearch(req models.SearchRequest, data models.Server) (bool, error) {
	fieldVal, err := getFieldValue(data, req.Field)
	if err != nil {
		return false, err
	}
	return !reflect.DeepEqual(fieldVal, req.Query.Value), nil
}

func (s *searchService) MatchSearch(req models.SearchRequest, data models.Server) (bool, error) {
	fieldVal, err := getFieldValue(data, req.Field)
	if err != nil {
		return false, err
	}
	strVal, ok := fieldVal.(string)
	if !ok {
		return false, fmt.Errorf("cannot apply match operator to non-string field: %s", req.Field)
	}
	return strings.Contains(strings.ToLower(strVal), strings.ToLower(req.Query.Value.(string))), nil
}

func (s *searchService) InSearch(req models.SearchRequest, data models.Server) (bool, error) {
	fieldVal, err := getFieldValue(data, req.Field)
	if err != nil {
		return false, err
	}
	arrVal, ok := fieldVal.([]interface{})
	if !ok {
		return false, fmt.Errorf("cannot apply in operator to non-array field: %s", req.Field)
	}
	for _, item := range arrVal {
		if reflect.DeepEqual(item, req.Query.Value) {
			return true, nil
		}
	}
	return false, nil
}

func (s *searchService) GteSearch(req models.SearchRequest, data models.Server) (bool, error) {
	fieldVal, err := getFieldValue(data, req.Field)
	if err != nil {
		return false, err
	}
	neededInt, err := castToInt(req.Query.Value)
	if err != nil {
		return false, err
	}
	serverInt, err := castToInt(fieldVal)
	if err != nil {
		return false, err
	}
	return serverInt >= neededInt, nil
}

func (s *searchService) GtSearch(req models.SearchRequest, data models.Server) (bool, error) {
	fieldVal, err := getFieldValue(data, req.Field)
	if err != nil {
		return false, err
	}
	neededInt, err := castToInt(req.Query.Value)
	if err != nil {
		return false, err
	}
	serverInt, err := castToInt(fieldVal)
	if err != nil {
		return false, err
	}
	return serverInt > neededInt, nil
}

func (s *searchService) LteSearch(req models.SearchRequest, data models.Server) (bool, error) {
	fieldVal, err := getFieldValue(data, req.Field)
	if err != nil {
		return false, err
	}
	neededInt, err := castToInt(req.Query.Value)
	if err != nil {
		return false, err
	}
	serverInt, err := castToInt(fieldVal)
	if err != nil {
		return false, err
	}
	return serverInt <= neededInt, nil
}

func (s *searchService) LtSearch(req models.SearchRequest, data models.Server) (bool, error) {
	fieldVal, err := getFieldValue(data, req.Field)
	if err != nil {
		return false, err
	}
	neededInt, err := castToInt(req.Query.Value)
	if err != nil {
		return false, err
	}
	serverInt, err := castToInt(fieldVal)
	if err != nil {
		return false, nil
	}
	return serverInt < neededInt, nil
}

func getFieldValue(data models.Server, field string) (interface{}, error) {
	switch field {
	case "name":
		return data.Name, nil
	case "capacity":
		return data.Capacity, nil
	case "current_player_count":
		return data.CurrentPlayerCount, nil
	default:
		val, ok := data.CustomData[field]
		if !ok {
			return nil, fmt.Errorf("field not found: %s", field)
		}
		return val, nil
	}
}

func castToInt(data interface{}) (int, error) {
	myType := reflect.TypeOf(data)
	if k := myType.Kind(); k == reflect.Int {
		return data.(int), nil
	} else if k == reflect.Float64 {
		return int(data.(float64)), nil
	}
	return 0, errors.New("invalid type")
}
