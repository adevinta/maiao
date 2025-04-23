class Maiao < Formula
  desc "Seamless GitHub PR management from the command-line"
  homepage "https://github.com/adevinta/maiao"
  url "https://github.com/adevinta/maiao.git",
    tag:      "1.2.1",
    revision: "8eba283e6770e33fc845aed370b9f0e121ec432d"
  license "MIT"
  conflicts_with "git-review"
  head "https://github.com/adevinta/maiao.git",
    branch: "main"

  depends_on "go" => :build

  def install
    ldflags = %W[
      -s -w
      -X github.com/adevinta/maiao/pkg/version.Version=#{version}+homebrew-adevinta-maiao
    ]

    system "go", "build", *std_go_args(ldflags: ldflags, output: bin/"git-review"), "./cmd/maiao"
    generate_completions_from_executable(bin/"git-review", "completion")
  end

  test do
    assert_match "#{version}+homebrew-adevinta-maiao", shell_output("#{bin}/git-review version")
  end
end
