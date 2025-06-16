class Morama < Formula
  desc "A CLI tool for managing your watched movies and dramas"
  homepage "https://github.com/kiku99/morama"
  url "https://github.com/kiku99/morama/archive/refs/tags/v1.0.5.tar.gz"
  sha256 "d244c7120de20e3c737b048ef470fd312d4b83eaf2b976ec8b4dc68035640e5e"
  version "v1.0.5"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "."
  end

  test do
    system "#{bin}/morama", "version"
  end
end 