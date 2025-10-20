resource "aws_opensearchserverless_collection" "vec" {
  name = var.collection_name
  type = "VECTORSEARCH"

  depends_on = [
    aws_opensearchserverless_security_policy.encryption,
    aws_opensearchserverless_security_policy.network
  ]
}

resource "aws_opensearchserverless_security_policy" "encryption" {
  name = "knowledge-oss-encryption"
  type = "encryption"

  policy = jsonencode({
    Rules = [
      {
        ResourceType = "collection"
        Resource     = ["collection/${var.collection_name}"]
      }
    ]
    AWSOwnedKey = true
  })
}

resource "aws_opensearchserverless_security_policy" "network" {
  name = "knowledge-oss-network"
  type = "network"
  policy = jsonencode([
    {
      Description = "PoC: allow public access to the collection"
      Rules = [
        {
          ResourceType = "collection"
          Resource     = ["collection/${var.collection_name}"]
        }
      ]
      AllowFromPublic = var.allow_public
    }
  ])
}

resource "aws_opensearchserverless_access_policy" "access_policy" {
  name        = "knowledge-oss-access"
  type        = "data"
  description = "Access control for OpenSearch Serverless"

  policy = jsonencode([
    {
      Description = "PoC: principals can manage and read/write documents"
      Principal   = var.principals
      Rules = [
        {
          ResourceType = "collection"
          Resource     = ["collection/${var.collection_name}"]
          Permission = [
            "aoss:CreateIndex",
            "aoss:UpdateIndex",
            "aoss:DescribeIndex",
            "aoss:ReadDocument",
            "aoss:WriteDocument"
          ]
        }
      ]
    }
  ])
}
