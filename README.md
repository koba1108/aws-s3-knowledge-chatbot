# aws-s3-knowledge-chatbot
aws-s3-knowledge-chatbot

## OpenSearch Index Setup

```
PUT /knowledge-base-index
{
  "settings": {
    "index": {
      "knn": true,
      "knn.algo_param.ef_search": 512
    }
  },
  "mappings": {
    "properties": {
      "vector": {
        "type": "knn_vector",
        "dimension": 1024,
        "method": {
          "name": "hnsw",
          "engine": "faiss",
          "parameters": {},
          "space_type": "l2"
        }
      },
      "text": {
        "type": "text",
        "index": true
      },
      "metadata": {
        "type": "object"
      }
    }
  }
}
```
