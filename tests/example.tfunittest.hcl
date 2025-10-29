test "vpc cidr validation" {
  module = "../examples/vpc"

  vars = {
    cidr_block = "10.0.0.0/16"
    mock_mode  = true
  }

  assert "cidr matches variable" {
    actual    = "output.vpc_cidr"
    expected  = "var.cidr_block"
    condition = "equals"
  }
}
