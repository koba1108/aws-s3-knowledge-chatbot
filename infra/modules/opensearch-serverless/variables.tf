variable "collection_name" {
  description = "Name of the OpenSearch Serverless collection"
  type        = string
}

variable "allow_public" {
  description = "Allow public access for PoC"
  type        = bool
  default     = false
}

variable "principals" {
  description = "List of IAM principals allowed to access the collection"
  type        = list(string)
  validation {
    condition     = length(var.principals) > 0
    error_message = "You must set at least one IAM principal ARN in var.principals."
  }
}
