package util

import (
	"errors"
	"fmt"
	"strings"

	e "github.com/olivere/elastic/v7"
)

var (
	allowedBoolClauses = []string{"must", "filter", "must_not", "should"}
)

type ElasticQuery struct {
	clause ElasticClause
	// Scoring define whether we'll use filter to query data from elasticsearch
	// since there is no scoring in filter, the search and cache in elasticsearch
	// can be faster and cheaper
	Scoring bool
}

// Source is used to implement olivere/elastic.v5's Query interface
func (e ElasticQuery) Source() (interface{}, error) {
	content, err := e.clause.Source()
	if err != nil {
		return nil, err
	}

	if e.Scoring {
		return map[string]interface{}{
			e.clause.Type(): content,
		}, nil
	}

	return map[string]interface{}{
		"bool": map[string]interface{}{
			"filter": map[string]interface{}{
				e.clause.Type(): content,
			},
		},
	}, nil
}

func (e *ElasticQuery) SetClause(clause ElasticClause) {
	e.clause = clause
}

func (e *ElasticQuery) SwitchScore(b bool) {
	e.Scoring = b
}

type ElasticClause interface {
	// Type should return the formal name for each clase's type,
	// like term, terms, bool, range etc.
	Type() string
	// Source returns the JSON-serializable result
	Source() (interface{}, error)
}

type MatchPhraseCause struct {
	key   string
	value string
}

func (MatchPhraseCause) Type() string {
	return "match_phrase"
}

func (m MatchPhraseCause) Source() (interface{}, error) {
	if m.key == "" {
		return nil, errors.New("Can not get source for empty MatchPhraseCause")
	}

	return map[string]interface{}{
		m.key: m.value,
	}, nil
}

func (m *MatchPhraseCause) Set(key, value string) {
	m.SetKey(key)
	m.SetValue(value)
}

func (m *MatchPhraseCause) SetKey(key string) {
	m.key = key
}

func (m *MatchPhraseCause) SetValue(value string) {
	m.value = value
}

/*
 *	Clauses are the implementation of each elastic clause
 *  Add more if needed
 */
type MatchClause struct {
	key                string
	value              string
	operator           string
	zeroTermsQuery     string
	minimumShouldMatch int
	fuzzyOn            bool
	fuzziness          string
	prefixLength       int
	maxExpansions      int
}

func (MatchClause) Type() string {
	return "match"
}

func (m MatchClause) Source() (interface{}, error) {
	if m.key == "" {
		return nil, errors.New("Can not get source for empty MatchClause")
	}

	if m.operator == "" {
		m.operator = "or"
	}

	if m.zeroTermsQuery == "" {
		m.zeroTermsQuery = "none"
	}

	if m.minimumShouldMatch == 0 {
		m.minimumShouldMatch = 1
	}

	if m.fuzziness == "" {
		m.fuzziness = "AUTO"
	}

	if m.prefixLength == 0 {
		m.prefixLength = 0
	}

	if m.maxExpansions == 0 {
		m.maxExpansions = 50
	}

	var result map[string]interface{}

	if !m.fuzzyOn {
		result = map[string]interface{}{
			m.key: map[string]interface{}{
				"query":                m.value,
				"operator":             m.operator,
				"zero_terms_query":     m.zeroTermsQuery,
				"minimum_should_match": m.minimumShouldMatch,
			},
		}
	} else {
		result = map[string]interface{}{
			m.key: map[string]interface{}{
				"query":                m.value,
				"operator":             m.operator,
				"zero_terms_query":     m.zeroTermsQuery,
				"minimum_should_match": m.minimumShouldMatch,
				"fuzziness":            m.fuzziness,
				"max_expansions":       m.maxExpansions,
				"prefix_length":        m.prefixLength,
			},
		}
	}

	return result, nil
}

func (m *MatchClause) Set(key, value string) {
	m.SetKey(key)
	m.SetValue(value)
}

func (m *MatchClause) SetKey(key string) {
	m.key = key
}

func (m *MatchClause) SetValue(value string) {
	m.value = value
}

func (m *MatchClause) SetOperator(operator string) {
	m.operator = operator
}

func (m *MatchClause) SetMinimumShouldMatch(number int) {
	m.minimumShouldMatch = number
}

func (m *MatchClause) SwitchFuzzy(b bool) {
	m.fuzzyOn = b
}

func (m *MatchClause) SetFuzziness(fuzziness string, prefix, expansion int) {
	m.fuzziness = fuzziness
	m.prefixLength = prefix
	m.maxExpansions = expansion
}

type TermClause struct {
	key      string
	value    interface{}
	boost    float64
	negative bool
}

func (t TermClause) Type() string {
	if t.negative {
		return "bool"
	} else {
		return "term"
	}
}

func (t TermClause) Source() (interface{}, error) {
	if t.key == "" {
		return nil, errors.New("Can not get source for empty TermClause")
	}

	if t.negative {
		boolClause := &BoolClause{}
		mustNotClause := &MustNotClause{}
		termClause := &TermClause{}
		boolClause.Set(mustNotClause)
		mustNotClause.Add(termClause)
		termClause.Set(t.key, t.value)
		termClause.SetBoost(t.boost)
		return boolClause.Source()
	}

	if t.boost > 0 {
		return map[string]interface{}{
			t.key: map[string]interface{}{
				"value": t.value,
				"boost": t.boost,
			},
		}, nil
	}

	return map[string]interface{}{
		t.key: t.value,
	}, nil
}

func (t *TermClause) Set(key string, value interface{}) {
	t.key = key
	t.value = value
}

func (t *TermClause) SetBoost(b float64) {
	t.boost = b
}

// 若 negative 为 true，那么 Source 生成的 query 会在 term 外面包一层 must_not
func (t *TermClause) SetNegative(n bool) {
	t.negative = n
}

type TermsClause struct {
	key    string
	values interface{}
}

func (TermsClause) Type() string {
	return "terms"
}

func (t TermsClause) Source() (interface{}, error) {
	if t.key == "" {
		return nil, errors.New("Can not get source for empty TermsClause")
	}

	return map[string]interface{}{
		t.key: t.values,
	}, nil
}

func (t *TermsClause) Set(key string, values interface{}) {
	t.SetKey(key)
	t.SetValue(values)
}

func (t *TermsClause) SetKey(key string) {
	t.key = key
}

func (t *TermsClause) SetValue(values interface{}) {
	t.values = values
}

type RangeClause struct {
	key     string
	content map[string]interface{}
}

func (RangeClause) Type() string {
	return "range"
}

func (r RangeClause) Source() (interface{}, error) {
	if r.key == "" {
		return nil, errors.New("Can not get source for empty RangeClause")
	}

	return map[string]interface{}{
		r.key: r.content,
	}, nil
}

func (r *RangeClause) SetKey(key string) {
	r.key = key
}

func (r *RangeClause) SetGT(value interface{}) {
	if r.content == nil {
		r.content = map[string]interface{}{}
	}

	r.content["gt"] = value
}

func (r *RangeClause) SetGTE(value interface{}) {
	if r.content == nil {
		r.content = map[string]interface{}{}
	}

	r.content["gte"] = value
}

func (r *RangeClause) SetLT(value interface{}) {
	if r.content == nil {
		r.content = map[string]interface{}{}
	}

	r.content["lt"] = value
}

func (r *RangeClause) SetLTE(value interface{}) {
	if r.content == nil {
		r.content = map[string]interface{}{}
	}

	r.content["lte"] = value
}

type ExistsClause struct {
	key       string
	existance bool
}

func (e ExistsClause) Type() string {
	if e.existance {
		return "exists"
	} else {
		// we can only use the "must_not" clause to express
		// inexist fields, so if the existance == false, the type
		// shall be "bool"
		return "bool"
	}
}

func (e ExistsClause) Source() (interface{}, error) {
	if e.key == "" {
		return nil, errors.New("Can not get source for empty ExistsClause")
	}

	if e.existance {
		return map[string]string{
			"field": e.key,
		}, nil
	}

	boolClause := &BoolClause{}
	mustNotClause := &MustNotClause{}
	existClause := &ExistsClause{}

	boolClause.Set(mustNotClause)
	mustNotClause.Add(existClause)
	existClause.Set(e.key, true)

	return boolClause.Source()
}

func (e *ExistsClause) Set(key string, exists bool) {
	e.SetKey(key)
	e.SetExistance(exists)
}

func (e *ExistsClause) SetKey(key string) {
	e.key = key
}

func (e *ExistsClause) SetExistance(value bool) {
	e.existance = value
}

type NestedClause struct {
	path      string
	scoreMode string
	clause    ElasticClause
}

func (NestedClause) Type() string {
	return "nested"
}

func (n NestedClause) Source() (interface{}, error) {
	content, err := n.clause.Source()
	if err != nil || n.path == "" {
		return nil, errors.New("Can not get source of invalid nested clause")
	}

	if n.scoreMode == "" {
		n.scoreMode = "none"
	}

	return map[string]interface{}{
		"path":       n.path,
		"score_mode": n.scoreMode,
		"query": map[string]interface{}{
			n.clause.Type(): content,
		},
	}, nil
}

func (n *NestedClause) SetPath(path string) {
	n.path = path
}

func (n *NestedClause) SetClause(clause ElasticClause) {
	n.clause = clause
}

func (n *NestedClause) SetScoreMode(mode string) {
	n.scoreMode = mode
}

type BoolClause struct {
	MinimumShouldMatch int
	clauses            map[string]ElasticClause
}

func (b BoolClause) Type() string {
	return "bool"
}

func (b BoolClause) Source() (interface{}, error) {
	if len(b.clauses) == 0 {
		return nil, errors.New("Can not get source for empty BoolClause")
	}

	resultMap := map[string]interface{}{
		"minimum_should_match": b.MinimumShouldMatch,
	}

	for typeName, clause := range b.clauses {
		content, err := clause.Source()
		if err != nil {
			continue
		}

		resultMap[typeName] = content
	}

	return resultMap, nil
}

func (b *BoolClause) Set(clause ElasticClause) error {
	if !strContains(allowedBoolClauses, clause.Type()) {
		return errors.New("Can not add current type of ElasticClause")
	}

	if _, exists := b.clauses[clause.Type()]; exists {
		return errors.New("Already have same clause")
	}

	if b.clauses == nil {
		b.clauses = map[string]ElasticClause{}
	}

	if clause.Type() == "should" {
		b.MinimumShouldMatch = 1
	}

	b.clauses[clause.Type()] = clause

	return nil
}

func (b *BoolClause) Del(typeName string) {
	delete(b.clauses, typeName)
}

func strContains(arr []string, str string) bool {
	for _, value := range arr {
		if value == str {
			return true
		}
	}

	return false
}

type BaseBoolClause struct {
	clauses []ElasticClause
}

func (BaseBoolClause) Type() string {
	return "baseBool"
}

func (b BaseBoolClause) Source() (interface{}, error) {
	if len(b.clauses) == 0 {
		return nil, errors.New("Can not get source for empty clause")
	}

	result := []interface{}{}
	for _, clause := range b.clauses {
		content, err := clause.Source()
		if err != nil {
			continue
		}

		result = append(result, map[string]interface{}{
			clause.Type(): content,
		})
	}

	return result, nil
}

func (b *BaseBoolClause) Add(clause ElasticClause) {
	b.clauses = append(b.clauses, clause)
}

func (b *BaseBoolClause) AddMulti(clauses ...ElasticClause) {
	b.clauses = append(b.clauses, clauses...)
}

func (b BaseBoolClause) Get(index int) ElasticClause {
	if index >= len(b.clauses) {
		return nil
	}

	return b.clauses[index]
}

func (b *BaseBoolClause) Del(index int) {
	if index >= len(b.clauses) {
		return
	}

	b.clauses = append(b.clauses[:index], b.clauses[index+1:]...)
}

func (b *BaseBoolClause) Count() int {
	return len(b.clauses)
}

type MustClause struct {
	BaseBoolClause
}

func (MustClause) Type() string {
	return "must"
}

type FilterClause struct {
	BaseBoolClause
}

func (FilterClause) Type() string {
	return "filter"
}

type MustNotClause struct {
	BaseBoolClause
}

func (MustNotClause) Type() string {
	return "must_not"
}

type ShouldClause struct {
	BaseBoolClause
}

func (ShouldClause) Type() string {
	return "should"
}

type FuzzyClause struct {
	key           string
	value         string
	boost         float32
	fuzziness     int
	prefixLength  int
	maxExpansions int
}

func (FuzzyClause) Type() string {
	return "fuzzy"
}

func (f FuzzyClause) Source() (interface{}, error) {
	if f.key == "" || f.value == "" {
		return nil, errors.New("Can not get source for empty clause")
	}

	configs := map[string]interface{}{
		"value": f.value,
	}

	if f.boost != 0 {
		configs["boost"] = f.boost
	}

	if f.fuzziness != 0 {
		configs["fuzziness"] = f.fuzziness
	} else {
		configs["fuzziness"] = 2
	}

	if f.prefixLength != 0 {
		configs["prefix_length"] = f.prefixLength
	} else {
		configs["prefix_length"] = 0
	}

	if f.maxExpansions != 0 {
		configs["max_expansions"] = f.maxExpansions
	} else {
		configs["max_expansions"] = 50
	}

	return map[string]interface{}{
		f.key: configs,
	}, nil
}

func (f *FuzzyClause) Set(key, value string) {
	f.SetKey(key)
	f.SetValue(value)
}

func (f *FuzzyClause) SetKey(key string) {
	f.key = key
}

func (f *FuzzyClause) SetValue(value string) {
	f.value = value
}

func (f *FuzzyClause) SetOptions(boost float32, fuzziness, prefixLength, maxExpansions int) {
	f.boost = boost
	f.fuzziness = fuzziness
	f.prefixLength = prefixLength
	f.maxExpansions = maxExpansions
}

type QueryStringClause struct {
	fields     []string
	querys     []string
	concatType string
}

func (QueryStringClause) Type() string {
	return "query_string"
}

func (q QueryStringClause) Source() (interface{}, error) {
	concat := q.concatType
	if concat == "" {
		concat = "AND"
	}
	queryString := strings.Join(q.querys, fmt.Sprintf(" %s ", concat))
	queryStringQuery := e.NewQueryStringQuery(queryString)

	for _, field := range q.fields {
		queryStringQuery.Field(field)
	}

	// the result is complete {"query_string": {...}},
	// but we only need the {...} part
	result, err := queryStringQuery.Source()
	if err != nil {
		return map[string]interface{}{}, err
	}

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		return map[string]interface{}{}, errors.New("is not query_string clause")
	}

	return resultMap[q.Type()], nil
}

func (q *QueryStringClause) AddQuery(query string) *QueryStringClause {
	if q.querys == nil {
		q.querys = make([]string, 0)
	}
	q.querys = append(q.querys, query)
	return q
}

func (q *QueryStringClause) AddQuerys(querys []string) *QueryStringClause {
	if q.querys == nil {
		q.querys = querys
		return q
	}
	q.querys = append(q.querys, querys...)
	return q
}

func (q *QueryStringClause) AddField(field string) *QueryStringClause {
	if q.fields == nil {
		q.fields = make([]string, 0)
	}
	q.fields = append(q.fields, field)
	return q
}

func (q *QueryStringClause) SetAnd() *QueryStringClause {
	q.concatType = "AND"
	return q
}

func (q *QueryStringClause) SetOr() *QueryStringClause {
	q.concatType = "OR"
	return q
}

type WildcardClause struct {
	name     string
	wildcard string
}

func (w *WildcardClause) Set(name, wildcard string) {
	w.name = name
	w.wildcard = wildcard
}

func (WildcardClause) Type() string {
	return "wildcard"
}

func (w WildcardClause) Source() (interface{}, error) {
	if w.name == "" || w.wildcard == "" {
		return nil, errors.New("Can not get source for empty WildcardClause")
	}

	return map[string]interface{}{
		w.name: map[string]interface{}{
			"value": w.wildcard,
		},
	}, nil
}

// ContainsClause 并不是 elasticsearch 中的 clause，
// 在 elasticsearch 5.2 版本中为了能实现 $in 或者 contains 功能，
// 需要使用 "constant_score": {"filter": {"terms": {"field": ["value"]}}} 这种 clause
// 为了能不用每次都写一大堆代码，我们把这种特殊又常用的例子封装在 ContainsClause 中
type ContainsClause struct {
	field  string
	values []string
}

func (ContainsClause) Type() string {
	return "constant_score" // 参考 ContainsClause 注释，它本质上是个 contant_score
}

func (c ContainsClause) Source() (interface{}, error) {
	if c.field == "" || len(c.values) == 0 {
		return nil, errors.New("Can not get source for empty ContainsClause")
	}

	termsClause := TermsClause{}
	termsClause.Set(c.field, c.values)
	source, _ := termsClause.Source()

	return map[string]interface{}{
		"filter": map[string]interface{}{
			termsClause.Type(): source,
		},
	}, nil
}

func (c *ContainsClause) Set(field string, values []string) {
	c.SetField(field)
	c.SetValues(values)
}

func (c *ContainsClause) SetField(field string) {
	c.field = field
}

func (c *ContainsClause) SetValues(values []string) {
	c.values = values
}
