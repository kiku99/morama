class Morama < Formula
  desc "Morama CLI tool"
  homepage "https://github.com/kiku99/morama"
  url "https://github.com/kiku99/morama/archive/refs/tags/v1.0.0.tar.gz"
  sha256 "YOUR_TARBALL_SHA256"
  version "v1.0.0"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "./cmd/morama"
  end

  test do
    system "#{bin}/morama", "--version"
  end
end 