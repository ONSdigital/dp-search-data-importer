// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"context"
	dpelasticsearch "github.com/ONSdigital/dp-elasticsearch/v3/client"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-search-data-importer/service"
	"sync"
)

// Ensure, that ElasticSearchMock does implement service.ElasticSearch.
// If this is not the case, regenerate this file with moq.
var _ service.ElasticSearch = &ElasticSearchMock{}

// ElasticSearchMock is a mock implementation of service.ElasticSearch.
//
//	func TestSomethingThatUsesElasticSearch(t *testing.T) {
//
//		// make and configure a mocked service.ElasticSearch
//		mockedElasticSearch := &ElasticSearchMock{
//			AddDocumentFunc: func(ctx context.Context, indexName string, documentID string, document []byte, opts *dpelasticsearch.AddDocumentOptions) error {
//				panic("mock out the AddDocument method")
//			},
//			BulkIndexAddFunc: func(ctx context.Context, action dpelasticsearch.BulkIndexerAction, index string, documentID string, document []byte) error {
//				panic("mock out the BulkIndexAdd method")
//			},
//			BulkIndexCloseFunc: func(contextMoqParam context.Context) error {
//				panic("mock out the BulkIndexClose method")
//			},
//			BulkUpdateFunc: func(ctx context.Context, indexName string, url string, settings []byte) ([]byte, error) {
//				panic("mock out the BulkUpdate method")
//			},
//			CheckerFunc: func(ctx context.Context, state *healthcheck.CheckState) error {
//				panic("mock out the Checker method")
//			},
//			CountIndicesFunc: func(ctx context.Context, indices []string) ([]byte, error) {
//				panic("mock out the CountIndices method")
//			},
//			CreateIndexFunc: func(ctx context.Context, indexName string, indexSettings []byte) error {
//				panic("mock out the CreateIndex method")
//			},
//			DeleteIndexFunc: func(ctx context.Context, indexName string) error {
//				panic("mock out the DeleteIndex method")
//			},
//			DeleteIndicesFunc: func(ctx context.Context, indices []string) error {
//				panic("mock out the DeleteIndices method")
//			},
//			GetAliasFunc: func(ctx context.Context) ([]byte, error) {
//				panic("mock out the GetAlias method")
//			},
//			GetIndicesFunc: func(ctx context.Context, indexPatterns []string) ([]byte, error) {
//				panic("mock out the GetIndices method")
//			},
//			MultiSearchFunc: func(ctx context.Context, searches []dpelasticsearch.Search) ([]byte, error) {
//				panic("mock out the MultiSearch method")
//			},
//			NewBulkIndexerFunc: func(contextMoqParam context.Context) error {
//				panic("mock out the NewBulkIndexer method")
//			},
//			SearchFunc: func(ctx context.Context, search dpelasticsearch.Search) ([]byte, error) {
//				panic("mock out the Search method")
//			},
//			UpdateAliasesFunc: func(ctx context.Context, alias string, removeIndices []string, addIndices []string) error {
//				panic("mock out the UpdateAliases method")
//			},
//		}
//
//		// use mockedElasticSearch in code that requires service.ElasticSearch
//		// and then make assertions.
//
//	}
type ElasticSearchMock struct {
	// AddDocumentFunc mocks the AddDocument method.
	AddDocumentFunc func(ctx context.Context, indexName string, documentID string, document []byte, opts *dpelasticsearch.AddDocumentOptions) error

	// BulkIndexAddFunc mocks the BulkIndexAdd method.
	BulkIndexAddFunc func(ctx context.Context, action dpelasticsearch.BulkIndexerAction, index string, documentID string, document []byte) error

	// BulkIndexCloseFunc mocks the BulkIndexClose method.
	BulkIndexCloseFunc func(contextMoqParam context.Context) error

	// BulkUpdateFunc mocks the BulkUpdate method.
	BulkUpdateFunc func(ctx context.Context, indexName string, url string, settings []byte) ([]byte, error)

	// CheckerFunc mocks the Checker method.
	CheckerFunc func(ctx context.Context, state *healthcheck.CheckState) error

	// CountIndicesFunc mocks the CountIndices method.
	CountIndicesFunc func(ctx context.Context, indices []string) ([]byte, error)

	// CreateIndexFunc mocks the CreateIndex method.
	CreateIndexFunc func(ctx context.Context, indexName string, indexSettings []byte) error

	// DeleteIndexFunc mocks the DeleteIndex method.
	DeleteIndexFunc func(ctx context.Context, indexName string) error

	// DeleteIndicesFunc mocks the DeleteIndices method.
	DeleteIndicesFunc func(ctx context.Context, indices []string) error

	// GetAliasFunc mocks the GetAlias method.
	GetAliasFunc func(ctx context.Context) ([]byte, error)

	// GetIndicesFunc mocks the GetIndices method.
	GetIndicesFunc func(ctx context.Context, indexPatterns []string) ([]byte, error)

	// MultiSearchFunc mocks the MultiSearch method.
	MultiSearchFunc func(ctx context.Context, searches []dpelasticsearch.Search) ([]byte, error)

	// NewBulkIndexerFunc mocks the NewBulkIndexer method.
	NewBulkIndexerFunc func(contextMoqParam context.Context) error

	// SearchFunc mocks the Search method.
	SearchFunc func(ctx context.Context, search dpelasticsearch.Search) ([]byte, error)

	// UpdateAliasesFunc mocks the UpdateAliases method.
	UpdateAliasesFunc func(ctx context.Context, alias string, removeIndices []string, addIndices []string) error

	// calls tracks calls to the methods.
	calls struct {
		// AddDocument holds details about calls to the AddDocument method.
		AddDocument []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// IndexName is the indexName argument value.
			IndexName string
			// DocumentID is the documentID argument value.
			DocumentID string
			// Document is the document argument value.
			Document []byte
			// Opts is the opts argument value.
			Opts *dpelasticsearch.AddDocumentOptions
		}
		// BulkIndexAdd holds details about calls to the BulkIndexAdd method.
		BulkIndexAdd []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Action is the action argument value.
			Action dpelasticsearch.BulkIndexerAction
			// Index is the index argument value.
			Index string
			// DocumentID is the documentID argument value.
			DocumentID string
			// Document is the document argument value.
			Document []byte
		}
		// BulkIndexClose holds details about calls to the BulkIndexClose method.
		BulkIndexClose []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
		}
		// BulkUpdate holds details about calls to the BulkUpdate method.
		BulkUpdate []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// IndexName is the indexName argument value.
			IndexName string
			// URL is the url argument value.
			URL string
			// Settings is the settings argument value.
			Settings []byte
		}
		// Checker holds details about calls to the Checker method.
		Checker []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// State is the state argument value.
			State *healthcheck.CheckState
		}
		// CountIndices holds details about calls to the CountIndices method.
		CountIndices []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Indices is the indices argument value.
			Indices []string
		}
		// CreateIndex holds details about calls to the CreateIndex method.
		CreateIndex []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// IndexName is the indexName argument value.
			IndexName string
			// IndexSettings is the indexSettings argument value.
			IndexSettings []byte
		}
		// DeleteIndex holds details about calls to the DeleteIndex method.
		DeleteIndex []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// IndexName is the indexName argument value.
			IndexName string
		}
		// DeleteIndices holds details about calls to the DeleteIndices method.
		DeleteIndices []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Indices is the indices argument value.
			Indices []string
		}
		// GetAlias holds details about calls to the GetAlias method.
		GetAlias []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
		}
		// GetIndices holds details about calls to the GetIndices method.
		GetIndices []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// IndexPatterns is the indexPatterns argument value.
			IndexPatterns []string
		}
		// MultiSearch holds details about calls to the MultiSearch method.
		MultiSearch []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Searches is the searches argument value.
			Searches []dpelasticsearch.Search
		}
		// NewBulkIndexer holds details about calls to the NewBulkIndexer method.
		NewBulkIndexer []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
		}
		// Search holds details about calls to the Search method.
		Search []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Search is the search argument value.
			Search dpelasticsearch.Search
		}
		// UpdateAliases holds details about calls to the UpdateAliases method.
		UpdateAliases []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Alias is the alias argument value.
			Alias string
			// RemoveIndices is the removeIndices argument value.
			RemoveIndices []string
			// AddIndices is the addIndices argument value.
			AddIndices []string
		}
	}
	lockAddDocument    sync.RWMutex
	lockBulkIndexAdd   sync.RWMutex
	lockBulkIndexClose sync.RWMutex
	lockBulkUpdate     sync.RWMutex
	lockChecker        sync.RWMutex
	lockCountIndices   sync.RWMutex
	lockCreateIndex    sync.RWMutex
	lockDeleteIndex    sync.RWMutex
	lockDeleteIndices  sync.RWMutex
	lockGetAlias       sync.RWMutex
	lockGetIndices     sync.RWMutex
	lockMultiSearch    sync.RWMutex
	lockNewBulkIndexer sync.RWMutex
	lockSearch         sync.RWMutex
	lockUpdateAliases  sync.RWMutex
}

// AddDocument calls AddDocumentFunc.
func (mock *ElasticSearchMock) AddDocument(ctx context.Context, indexName string, documentID string, document []byte, opts *dpelasticsearch.AddDocumentOptions) error {
	if mock.AddDocumentFunc == nil {
		panic("ElasticSearchMock.AddDocumentFunc: method is nil but ElasticSearch.AddDocument was just called")
	}
	callInfo := struct {
		Ctx        context.Context
		IndexName  string
		DocumentID string
		Document   []byte
		Opts       *dpelasticsearch.AddDocumentOptions
	}{
		Ctx:        ctx,
		IndexName:  indexName,
		DocumentID: documentID,
		Document:   document,
		Opts:       opts,
	}
	mock.lockAddDocument.Lock()
	mock.calls.AddDocument = append(mock.calls.AddDocument, callInfo)
	mock.lockAddDocument.Unlock()
	return mock.AddDocumentFunc(ctx, indexName, documentID, document, opts)
}

// AddDocumentCalls gets all the calls that were made to AddDocument.
// Check the length with:
//
//	len(mockedElasticSearch.AddDocumentCalls())
func (mock *ElasticSearchMock) AddDocumentCalls() []struct {
	Ctx        context.Context
	IndexName  string
	DocumentID string
	Document   []byte
	Opts       *dpelasticsearch.AddDocumentOptions
} {
	var calls []struct {
		Ctx        context.Context
		IndexName  string
		DocumentID string
		Document   []byte
		Opts       *dpelasticsearch.AddDocumentOptions
	}
	mock.lockAddDocument.RLock()
	calls = mock.calls.AddDocument
	mock.lockAddDocument.RUnlock()
	return calls
}

// BulkIndexAdd calls BulkIndexAddFunc.
func (mock *ElasticSearchMock) BulkIndexAdd(ctx context.Context, action dpelasticsearch.BulkIndexerAction, index string, documentID string, document []byte) error {
	if mock.BulkIndexAddFunc == nil {
		panic("ElasticSearchMock.BulkIndexAddFunc: method is nil but ElasticSearch.BulkIndexAdd was just called")
	}
	callInfo := struct {
		Ctx        context.Context
		Action     dpelasticsearch.BulkIndexerAction
		Index      string
		DocumentID string
		Document   []byte
	}{
		Ctx:        ctx,
		Action:     action,
		Index:      index,
		DocumentID: documentID,
		Document:   document,
	}
	mock.lockBulkIndexAdd.Lock()
	mock.calls.BulkIndexAdd = append(mock.calls.BulkIndexAdd, callInfo)
	mock.lockBulkIndexAdd.Unlock()
	return mock.BulkIndexAddFunc(ctx, action, index, documentID, document)
}

// BulkIndexAddCalls gets all the calls that were made to BulkIndexAdd.
// Check the length with:
//
//	len(mockedElasticSearch.BulkIndexAddCalls())
func (mock *ElasticSearchMock) BulkIndexAddCalls() []struct {
	Ctx        context.Context
	Action     dpelasticsearch.BulkIndexerAction
	Index      string
	DocumentID string
	Document   []byte
} {
	var calls []struct {
		Ctx        context.Context
		Action     dpelasticsearch.BulkIndexerAction
		Index      string
		DocumentID string
		Document   []byte
	}
	mock.lockBulkIndexAdd.RLock()
	calls = mock.calls.BulkIndexAdd
	mock.lockBulkIndexAdd.RUnlock()
	return calls
}

// BulkIndexClose calls BulkIndexCloseFunc.
func (mock *ElasticSearchMock) BulkIndexClose(contextMoqParam context.Context) error {
	if mock.BulkIndexCloseFunc == nil {
		panic("ElasticSearchMock.BulkIndexCloseFunc: method is nil but ElasticSearch.BulkIndexClose was just called")
	}
	callInfo := struct {
		ContextMoqParam context.Context
	}{
		ContextMoqParam: contextMoqParam,
	}
	mock.lockBulkIndexClose.Lock()
	mock.calls.BulkIndexClose = append(mock.calls.BulkIndexClose, callInfo)
	mock.lockBulkIndexClose.Unlock()
	return mock.BulkIndexCloseFunc(contextMoqParam)
}

// BulkIndexCloseCalls gets all the calls that were made to BulkIndexClose.
// Check the length with:
//
//	len(mockedElasticSearch.BulkIndexCloseCalls())
func (mock *ElasticSearchMock) BulkIndexCloseCalls() []struct {
	ContextMoqParam context.Context
} {
	var calls []struct {
		ContextMoqParam context.Context
	}
	mock.lockBulkIndexClose.RLock()
	calls = mock.calls.BulkIndexClose
	mock.lockBulkIndexClose.RUnlock()
	return calls
}

// BulkUpdate calls BulkUpdateFunc.
func (mock *ElasticSearchMock) BulkUpdate(ctx context.Context, indexName string, url string, settings []byte) ([]byte, error) {
	if mock.BulkUpdateFunc == nil {
		panic("ElasticSearchMock.BulkUpdateFunc: method is nil but ElasticSearch.BulkUpdate was just called")
	}
	callInfo := struct {
		Ctx       context.Context
		IndexName string
		URL       string
		Settings  []byte
	}{
		Ctx:       ctx,
		IndexName: indexName,
		URL:       url,
		Settings:  settings,
	}
	mock.lockBulkUpdate.Lock()
	mock.calls.BulkUpdate = append(mock.calls.BulkUpdate, callInfo)
	mock.lockBulkUpdate.Unlock()
	return mock.BulkUpdateFunc(ctx, indexName, url, settings)
}

// BulkUpdateCalls gets all the calls that were made to BulkUpdate.
// Check the length with:
//
//	len(mockedElasticSearch.BulkUpdateCalls())
func (mock *ElasticSearchMock) BulkUpdateCalls() []struct {
	Ctx       context.Context
	IndexName string
	URL       string
	Settings  []byte
} {
	var calls []struct {
		Ctx       context.Context
		IndexName string
		URL       string
		Settings  []byte
	}
	mock.lockBulkUpdate.RLock()
	calls = mock.calls.BulkUpdate
	mock.lockBulkUpdate.RUnlock()
	return calls
}

// Checker calls CheckerFunc.
func (mock *ElasticSearchMock) Checker(ctx context.Context, state *healthcheck.CheckState) error {
	if mock.CheckerFunc == nil {
		panic("ElasticSearchMock.CheckerFunc: method is nil but ElasticSearch.Checker was just called")
	}
	callInfo := struct {
		Ctx   context.Context
		State *healthcheck.CheckState
	}{
		Ctx:   ctx,
		State: state,
	}
	mock.lockChecker.Lock()
	mock.calls.Checker = append(mock.calls.Checker, callInfo)
	mock.lockChecker.Unlock()
	return mock.CheckerFunc(ctx, state)
}

// CheckerCalls gets all the calls that were made to Checker.
// Check the length with:
//
//	len(mockedElasticSearch.CheckerCalls())
func (mock *ElasticSearchMock) CheckerCalls() []struct {
	Ctx   context.Context
	State *healthcheck.CheckState
} {
	var calls []struct {
		Ctx   context.Context
		State *healthcheck.CheckState
	}
	mock.lockChecker.RLock()
	calls = mock.calls.Checker
	mock.lockChecker.RUnlock()
	return calls
}

// CountIndices calls CountIndicesFunc.
func (mock *ElasticSearchMock) CountIndices(ctx context.Context, indices []string) ([]byte, error) {
	if mock.CountIndicesFunc == nil {
		panic("ElasticSearchMock.CountIndicesFunc: method is nil but ElasticSearch.CountIndices was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		Indices []string
	}{
		Ctx:     ctx,
		Indices: indices,
	}
	mock.lockCountIndices.Lock()
	mock.calls.CountIndices = append(mock.calls.CountIndices, callInfo)
	mock.lockCountIndices.Unlock()
	return mock.CountIndicesFunc(ctx, indices)
}

// CountIndicesCalls gets all the calls that were made to CountIndices.
// Check the length with:
//
//	len(mockedElasticSearch.CountIndicesCalls())
func (mock *ElasticSearchMock) CountIndicesCalls() []struct {
	Ctx     context.Context
	Indices []string
} {
	var calls []struct {
		Ctx     context.Context
		Indices []string
	}
	mock.lockCountIndices.RLock()
	calls = mock.calls.CountIndices
	mock.lockCountIndices.RUnlock()
	return calls
}

// CreateIndex calls CreateIndexFunc.
func (mock *ElasticSearchMock) CreateIndex(ctx context.Context, indexName string, indexSettings []byte) error {
	if mock.CreateIndexFunc == nil {
		panic("ElasticSearchMock.CreateIndexFunc: method is nil but ElasticSearch.CreateIndex was just called")
	}
	callInfo := struct {
		Ctx           context.Context
		IndexName     string
		IndexSettings []byte
	}{
		Ctx:           ctx,
		IndexName:     indexName,
		IndexSettings: indexSettings,
	}
	mock.lockCreateIndex.Lock()
	mock.calls.CreateIndex = append(mock.calls.CreateIndex, callInfo)
	mock.lockCreateIndex.Unlock()
	return mock.CreateIndexFunc(ctx, indexName, indexSettings)
}

// CreateIndexCalls gets all the calls that were made to CreateIndex.
// Check the length with:
//
//	len(mockedElasticSearch.CreateIndexCalls())
func (mock *ElasticSearchMock) CreateIndexCalls() []struct {
	Ctx           context.Context
	IndexName     string
	IndexSettings []byte
} {
	var calls []struct {
		Ctx           context.Context
		IndexName     string
		IndexSettings []byte
	}
	mock.lockCreateIndex.RLock()
	calls = mock.calls.CreateIndex
	mock.lockCreateIndex.RUnlock()
	return calls
}

// DeleteIndex calls DeleteIndexFunc.
func (mock *ElasticSearchMock) DeleteIndex(ctx context.Context, indexName string) error {
	if mock.DeleteIndexFunc == nil {
		panic("ElasticSearchMock.DeleteIndexFunc: method is nil but ElasticSearch.DeleteIndex was just called")
	}
	callInfo := struct {
		Ctx       context.Context
		IndexName string
	}{
		Ctx:       ctx,
		IndexName: indexName,
	}
	mock.lockDeleteIndex.Lock()
	mock.calls.DeleteIndex = append(mock.calls.DeleteIndex, callInfo)
	mock.lockDeleteIndex.Unlock()
	return mock.DeleteIndexFunc(ctx, indexName)
}

// DeleteIndexCalls gets all the calls that were made to DeleteIndex.
// Check the length with:
//
//	len(mockedElasticSearch.DeleteIndexCalls())
func (mock *ElasticSearchMock) DeleteIndexCalls() []struct {
	Ctx       context.Context
	IndexName string
} {
	var calls []struct {
		Ctx       context.Context
		IndexName string
	}
	mock.lockDeleteIndex.RLock()
	calls = mock.calls.DeleteIndex
	mock.lockDeleteIndex.RUnlock()
	return calls
}

// DeleteIndices calls DeleteIndicesFunc.
func (mock *ElasticSearchMock) DeleteIndices(ctx context.Context, indices []string) error {
	if mock.DeleteIndicesFunc == nil {
		panic("ElasticSearchMock.DeleteIndicesFunc: method is nil but ElasticSearch.DeleteIndices was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		Indices []string
	}{
		Ctx:     ctx,
		Indices: indices,
	}
	mock.lockDeleteIndices.Lock()
	mock.calls.DeleteIndices = append(mock.calls.DeleteIndices, callInfo)
	mock.lockDeleteIndices.Unlock()
	return mock.DeleteIndicesFunc(ctx, indices)
}

// DeleteIndicesCalls gets all the calls that were made to DeleteIndices.
// Check the length with:
//
//	len(mockedElasticSearch.DeleteIndicesCalls())
func (mock *ElasticSearchMock) DeleteIndicesCalls() []struct {
	Ctx     context.Context
	Indices []string
} {
	var calls []struct {
		Ctx     context.Context
		Indices []string
	}
	mock.lockDeleteIndices.RLock()
	calls = mock.calls.DeleteIndices
	mock.lockDeleteIndices.RUnlock()
	return calls
}

// GetAlias calls GetAliasFunc.
func (mock *ElasticSearchMock) GetAlias(ctx context.Context) ([]byte, error) {
	if mock.GetAliasFunc == nil {
		panic("ElasticSearchMock.GetAliasFunc: method is nil but ElasticSearch.GetAlias was just called")
	}
	callInfo := struct {
		Ctx context.Context
	}{
		Ctx: ctx,
	}
	mock.lockGetAlias.Lock()
	mock.calls.GetAlias = append(mock.calls.GetAlias, callInfo)
	mock.lockGetAlias.Unlock()
	return mock.GetAliasFunc(ctx)
}

// GetAliasCalls gets all the calls that were made to GetAlias.
// Check the length with:
//
//	len(mockedElasticSearch.GetAliasCalls())
func (mock *ElasticSearchMock) GetAliasCalls() []struct {
	Ctx context.Context
} {
	var calls []struct {
		Ctx context.Context
	}
	mock.lockGetAlias.RLock()
	calls = mock.calls.GetAlias
	mock.lockGetAlias.RUnlock()
	return calls
}

// GetIndices calls GetIndicesFunc.
func (mock *ElasticSearchMock) GetIndices(ctx context.Context, indexPatterns []string) ([]byte, error) {
	if mock.GetIndicesFunc == nil {
		panic("ElasticSearchMock.GetIndicesFunc: method is nil but ElasticSearch.GetIndices was just called")
	}
	callInfo := struct {
		Ctx           context.Context
		IndexPatterns []string
	}{
		Ctx:           ctx,
		IndexPatterns: indexPatterns,
	}
	mock.lockGetIndices.Lock()
	mock.calls.GetIndices = append(mock.calls.GetIndices, callInfo)
	mock.lockGetIndices.Unlock()
	return mock.GetIndicesFunc(ctx, indexPatterns)
}

// GetIndicesCalls gets all the calls that were made to GetIndices.
// Check the length with:
//
//	len(mockedElasticSearch.GetIndicesCalls())
func (mock *ElasticSearchMock) GetIndicesCalls() []struct {
	Ctx           context.Context
	IndexPatterns []string
} {
	var calls []struct {
		Ctx           context.Context
		IndexPatterns []string
	}
	mock.lockGetIndices.RLock()
	calls = mock.calls.GetIndices
	mock.lockGetIndices.RUnlock()
	return calls
}

// MultiSearch calls MultiSearchFunc.
func (mock *ElasticSearchMock) MultiSearch(ctx context.Context, searches []dpelasticsearch.Search) ([]byte, error) {
	if mock.MultiSearchFunc == nil {
		panic("ElasticSearchMock.MultiSearchFunc: method is nil but ElasticSearch.MultiSearch was just called")
	}
	callInfo := struct {
		Ctx      context.Context
		Searches []dpelasticsearch.Search
	}{
		Ctx:      ctx,
		Searches: searches,
	}
	mock.lockMultiSearch.Lock()
	mock.calls.MultiSearch = append(mock.calls.MultiSearch, callInfo)
	mock.lockMultiSearch.Unlock()
	return mock.MultiSearchFunc(ctx, searches)
}

// MultiSearchCalls gets all the calls that were made to MultiSearch.
// Check the length with:
//
//	len(mockedElasticSearch.MultiSearchCalls())
func (mock *ElasticSearchMock) MultiSearchCalls() []struct {
	Ctx      context.Context
	Searches []dpelasticsearch.Search
} {
	var calls []struct {
		Ctx      context.Context
		Searches []dpelasticsearch.Search
	}
	mock.lockMultiSearch.RLock()
	calls = mock.calls.MultiSearch
	mock.lockMultiSearch.RUnlock()
	return calls
}

// NewBulkIndexer calls NewBulkIndexerFunc.
func (mock *ElasticSearchMock) NewBulkIndexer(contextMoqParam context.Context) error {
	if mock.NewBulkIndexerFunc == nil {
		panic("ElasticSearchMock.NewBulkIndexerFunc: method is nil but ElasticSearch.NewBulkIndexer was just called")
	}
	callInfo := struct {
		ContextMoqParam context.Context
	}{
		ContextMoqParam: contextMoqParam,
	}
	mock.lockNewBulkIndexer.Lock()
	mock.calls.NewBulkIndexer = append(mock.calls.NewBulkIndexer, callInfo)
	mock.lockNewBulkIndexer.Unlock()
	return mock.NewBulkIndexerFunc(contextMoqParam)
}

// NewBulkIndexerCalls gets all the calls that were made to NewBulkIndexer.
// Check the length with:
//
//	len(mockedElasticSearch.NewBulkIndexerCalls())
func (mock *ElasticSearchMock) NewBulkIndexerCalls() []struct {
	ContextMoqParam context.Context
} {
	var calls []struct {
		ContextMoqParam context.Context
	}
	mock.lockNewBulkIndexer.RLock()
	calls = mock.calls.NewBulkIndexer
	mock.lockNewBulkIndexer.RUnlock()
	return calls
}

// Search calls SearchFunc.
func (mock *ElasticSearchMock) Search(ctx context.Context, search dpelasticsearch.Search) ([]byte, error) {
	if mock.SearchFunc == nil {
		panic("ElasticSearchMock.SearchFunc: method is nil but ElasticSearch.Search was just called")
	}
	callInfo := struct {
		Ctx    context.Context
		Search dpelasticsearch.Search
	}{
		Ctx:    ctx,
		Search: search,
	}
	mock.lockSearch.Lock()
	mock.calls.Search = append(mock.calls.Search, callInfo)
	mock.lockSearch.Unlock()
	return mock.SearchFunc(ctx, search)
}

// SearchCalls gets all the calls that were made to Search.
// Check the length with:
//
//	len(mockedElasticSearch.SearchCalls())
func (mock *ElasticSearchMock) SearchCalls() []struct {
	Ctx    context.Context
	Search dpelasticsearch.Search
} {
	var calls []struct {
		Ctx    context.Context
		Search dpelasticsearch.Search
	}
	mock.lockSearch.RLock()
	calls = mock.calls.Search
	mock.lockSearch.RUnlock()
	return calls
}

// UpdateAliases calls UpdateAliasesFunc.
func (mock *ElasticSearchMock) UpdateAliases(ctx context.Context, alias string, removeIndices []string, addIndices []string) error {
	if mock.UpdateAliasesFunc == nil {
		panic("ElasticSearchMock.UpdateAliasesFunc: method is nil but ElasticSearch.UpdateAliases was just called")
	}
	callInfo := struct {
		Ctx           context.Context
		Alias         string
		RemoveIndices []string
		AddIndices    []string
	}{
		Ctx:           ctx,
		Alias:         alias,
		RemoveIndices: removeIndices,
		AddIndices:    addIndices,
	}
	mock.lockUpdateAliases.Lock()
	mock.calls.UpdateAliases = append(mock.calls.UpdateAliases, callInfo)
	mock.lockUpdateAliases.Unlock()
	return mock.UpdateAliasesFunc(ctx, alias, removeIndices, addIndices)
}

// UpdateAliasesCalls gets all the calls that were made to UpdateAliases.
// Check the length with:
//
//	len(mockedElasticSearch.UpdateAliasesCalls())
func (mock *ElasticSearchMock) UpdateAliasesCalls() []struct {
	Ctx           context.Context
	Alias         string
	RemoveIndices []string
	AddIndices    []string
} {
	var calls []struct {
		Ctx           context.Context
		Alias         string
		RemoveIndices []string
		AddIndices    []string
	}
	mock.lockUpdateAliases.RLock()
	calls = mock.calls.UpdateAliases
	mock.lockUpdateAliases.RUnlock()
	return calls
}
