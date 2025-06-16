class Morama < Formula
  desc "A CLI tool for managing your watched movies and dramas"
  homepage "https://github.com/kiku99/morama"
  url "https://github.com/kiku99/morama/archive/refs/tags/v1.0.6.tar.gz"
  sha256 "94ede474745582542ac75c92066a17a8ff3a45b6fb182849bea60d5a697ed7b1"
  version "v1.0.6"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "."
  end

  test do
    system "#{bin}/morama", "version"
  end
end 