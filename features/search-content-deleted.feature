Feature: Search content deleted from elasticsearch

  Scenario: Delete event is processed and query is sent to Elasticsearch
    Given elasticsearch is healthy
    And elasticsearch returns the following response for bulk update
    """
      {
        "took": 13,
        "errors": false,
        "items": []
      }
    """
    When the service starts
    And this "search-content-deleted" event is queued, to be consumed:
    """
      {
        "uri":"some_deleted_uri",
        "search_index":"ons"
      }
    """
    Then this bulk delete is sent to elasticsearch for index "ons"
    """
      {"delete":{"_id":"some_deleted_uri"}}
    """
