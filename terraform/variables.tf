variable "prefix" {
  type    = string
  default = "arar-"
}

variable "tags" {
  type = map(string)
  default = {
    "arar" = "test"
  }
}
