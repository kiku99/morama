class Morama < Formula
  desc "A CLI tool for managing your watched movies and dramas"
  homepage "https://github.com/kiku99/morama"
  url "https://github.com/kiku99/morama/archive/refs/tags/v1.0.0.tar.gz"
  sha256 "0000000000000000000000000000000000000000000000000000000000000000"
  version "v1.0.0"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "."
  end

  test do
    system "#{bin}/morama", "version"
  end
end 