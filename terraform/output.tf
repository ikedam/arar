output "AWS_ACCESS_KEY_ID" {
  value = aws_iam_access_key.testuser_key.id
}

output "AWS_SECRET_ACCESS_KEY" {
  value     = aws_iam_access_key.testuser_key.secret
  sensitive = true
}

output "MFA_SERIAL_NUMBER" {
  value = aws_iam_virtual_mfa_device.testuser_mfa.arn
}

output "MFA_SEED" {
  value = aws_iam_virtual_mfa_device.testuser_mfa.base_32_string_seed
}
