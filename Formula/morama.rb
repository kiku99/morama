class Morama < Formula
  desc "A CLI tool for managing your watched movies and dramas"
  homepage "https://github.com/kiku99/morama"
  url "https://github.com/kiku99/morama/archive/refs/tags/v1.1.0.tar.gz"
  sha256 "ad420b1a1e90cb443a750b14d230e323a0a891bd0992cc9fbc746ce7758a53a4"
  version "v1.1.0"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "."
  end

  test do
    system "#{bin}/morama", "version"
  end
end 