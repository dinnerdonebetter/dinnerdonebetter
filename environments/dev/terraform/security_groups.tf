resource "aws_security_group" "allow_ssh" {
  name        = "ssh"
  description = "Allow inbound SSH traffic from any IP"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
