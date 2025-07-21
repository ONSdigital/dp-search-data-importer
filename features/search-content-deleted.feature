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
    And this delete event is queued, to be consumed
    """
      {
        "URI":"some_deleted_uri",
        "SearchIndex":"ons"
      }
    """
    Then this delete request is sent to elasticsearch
    """
      /ons/_doc/some_deleted_uri
    """
