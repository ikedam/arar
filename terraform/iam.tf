resource "aws_iam_user" "testuser" {
  name = "${var.prefix}testuser"
  tags = var.tags
}

resource "aws_iam_access_key" "testuser_key" {
  user = aws_iam_user.testuser.name
}

resource "aws_iam_virtual_mfa_device" "testuser_mfa" {
  virtual_mfa_device_name = aws_iam_user.testuser.name
  tags                    = var.tags
}
