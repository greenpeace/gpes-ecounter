package main

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	debug = new(bool)
	os.Exit(m.Run())
}

func TestFileToString(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "input.txt")
	content := "one@example.com\ntwo@example.com\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	got := fileToString(path)
	if got != content {
		t.Fatalf("expected %q, got %q", content, got)
	}
}

func TestFileToStringMissingFile(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") == "1" {
		fileToString(filepath.Join(t.TempDir(), "does-not-exist.txt"))
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestFileToStringMissingFile")
	cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected non-zero exit, got output: %s", out)
	}
	if !strings.Contains(string(out), "ERROR: The file/path") {
		t.Fatalf("expected error message, got: %s", out)
	}
}

func TestIsDirectory(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "a.txt")
	if err := os.WriteFile(file, []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	if !isDirectory(dir) {
		t.Errorf("expected %q to be a directory", dir)
	}
	if isDirectory(file) {
		t.Errorf("expected %q to not be a directory", file)
	}
	if isDirectory(filepath.Join(dir, "missing")) {
		t.Errorf("expected missing path to not be a directory")
	}
}

func TestSupportedExtension(t *testing.T) {
	extensions := []struct {
		path string
		want bool
	}{
		{"a.txt", true},
		{"a.csv", true},
		{"a.tsv", true},
		{"a.xml", true},
		{"a.html", true},
		{"a.json", true},
		{"a.TXT", true},
		{"a.JsOn", true},
		{"a.log", false},
		{"a.sql", false},
		{"a.pdf", false},
		{"a.txt.bak", false},
		{"noextension", false},
		{"dir/a.csv", true},
	}
	for _, tc := range extensions {
		if got := supportedExtension(tc.path); got != tc.want {
			t.Errorf("supportedExtension(%q) = %v, want %v", tc.path, got, tc.want)
		}
	}
}

func TestFolderToString(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "sub")
	if err := os.Mkdir(sub, 0755); err != nil {
		t.Fatal(err)
	}
	files := map[string]string{
		"a.txt":            "first@example.com\n",
		"b.csv":            "second@example.com\n",
		"c.tsv":            "third@example.com\n",
		"d.xml":            "fourth@example.com\n",
		"e.html":           "fifth@example.com\n",
		"f.json":           "sixth@example.com\n",
		"sub/g.txt":        "seventh@example.com\n",
		"ignore.log":       "ignored@example.com\n",
		"sub/ignore.txtx":  "also-ignored@example.com\n",
		"nested/noext":     "no-extension@example.com\n",
	}
	for name, content := range files {
		path := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}
	got := folderToString(dir)
	for _, name := range []string{"a.txt", "b.csv", "c.tsv", "d.xml", "e.html", "f.json", "sub/g.txt"} {
		want := files[name]
		if !strings.Contains(got, want) {
			t.Errorf("expected content of %q (%q) to be scanned, got: %q", name, want, got)
		}
	}
	for _, name := range []string{"ignore.log", "sub/ignore.txtx", "nested/noext"} {
		notWant := files[name]
		if strings.Contains(got, notWant) {
			t.Errorf("expected content of %q (%q) to be skipped, got: %q", name, notWant, got)
		}
	}
}

func TestInputToStringFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "input.txt")
	content := "file@example.com\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	if got := inputToString(path); got != content {
		t.Fatalf("expected %q, got %q", content, got)
	}
}

func TestInputToStringFolder(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("folder@example.com\n"), 0644); err != nil {
		t.Fatal(err)
	}
	got := inputToString(dir)
	if got != "folder@example.com\n" {
		t.Fatalf("expected folder content, got %q", got)
	}
}

func TestOutputInsideFolder(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "sub")
	if err := os.Mkdir(sub, 0755); err != nil {
		t.Fatal(err)
	}
	outside := filepath.Join(t.TempDir(), "out.txt")
	extensions := []struct {
		name   string
		output string
		want   bool
	}{
		{"inside folder", filepath.Join(dir, "out.txt"), true},
		{"inside subfolder", filepath.Join(sub, "out.txt"), true},
		{"inside folder with nested path", filepath.Join(dir, "sub", "deeper", "out.txt"), true},
		{"outside folder", outside, false},
		{"parent of folder", filepath.Join(filepath.Dir(dir), "out.txt"), false},
		{"sibling folder", filepath.Join(filepath.Dir(dir), "other", "out.txt"), false},
	}
	for _, tc := range extensions {
		if got := outputInsideFolder(dir, tc.output); got != tc.want {
			t.Errorf("%s: outputInsideFolder(%q, %q) = %v, want %v", tc.name, dir, tc.output, got, tc.want)
		}
	}
}

func TestOutputInsideFolderRelativePaths(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "sub"), 0755); err != nil {
		t.Fatal(err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	rel, err := filepath.Rel(cwd, dir)
	if err != nil {
		t.Fatal(err)
	}
	if !outputInsideFolder(rel, filepath.Join(rel, "out.txt")) {
		t.Errorf("expected output inside folder to be detected with relative paths")
	}
	if outputInsideFolder(rel, filepath.Join(rel, "..", "out.txt")) {
		t.Errorf("expected output outside folder to not be detected with relative paths")
	}
}

func TestSearchInString(t *testing.T) {
	extensions := []struct {
		name       string
		total      string
		expression string
		want       []string
	}{
		{"emails", "test@example.com and foo.bar@baz.io", emailRegex, []string{"test@example.com", "foo.bar@baz.io"}},
		{"sha256", "abc 9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08 xyz", shaRegex, []string{"9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"}},
		{"urls", "see https://example.com/path?q=1 now", urlsRegex, []string{"https://example.com/path"}},
		{"dnis", "dni 12345678Z and 87654321X", dninieRegex, []string{"12345678Z", "87654321X"}},
	{"no matches", "nothing here", emailRegex, nil},
}
	for _, tc := range extensions {
		got := searchInString(tc.total, tc.expression)
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("%s: searchInString(%q, %q) = %v, want %v", tc.name, tc.total, tc.expression, got, tc.want)
		}
	}
}

func TestArrayToLowercase(t *testing.T) {
	got := arrayToLowercase([]string{"TEST@EXAMPLE.COM", "Foo@Bar.IO", ""})
	want := []string{"test@example.com", "foo@bar.io", ""}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("arrayToLowercase = %v, want %v", got, want)
	}
	if got := arrayToLowercase(nil); got != nil {
		t.Errorf("arrayToLowercase(nil) = %v, want nil", got)
	}
}

func TestArrayToUpercase(t *testing.T) {
	got := arrayToUpercase([]string{"test@example.com", "foo@bar.io"})
	want := []string{"TEST@EXAMPLE.COM", "FOO@BAR.IO"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("arrayToUpercase = %v, want %v", got, want)
	}
	if got := arrayToUpercase(nil); got != nil {
		t.Errorf("arrayToUpercase(nil) = %v, want nil", got)
	}
}

func TestUniquesInArray(t *testing.T) {
	got := uniquesInArray([]string{"a", "b", "a", "c", "b", "a"})
	assertSet(t, got, []string{"a", "b", "c"})
	if got := uniquesInArray([]string{}); len(got) != 0 {
		t.Errorf("uniquesInArray([]) = %v, want empty", got)
	}
}

func TestStringToSha256(t *testing.T) {
	// Known sha256 of the string "test".
	got := stringToSha256("test")
	want := "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"
	if got != want {
		t.Errorf("stringToSha256(test) = %q, want %q", got, want)
	}
}

func TestArrayToSha256(t *testing.T) {
	got := arrayToSha256([]string{"test", "foo"})
	want := []string{
		"9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
		"2c26b46b68ffc68ff99b453c1d30413413422d706483bfa0f98a5e886266e7ae",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("arrayToSha256 = %v, want %v", got, want)
	}
}

func TestStringToFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.txt")
	content := "result@example.com\n"
	stringToFile(path, content)
	dat, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(dat) != content {
		t.Errorf("stringToFile wrote %q, want %q", string(dat), content)
	}
}

// readLines reads a file and returns its non-empty lines.
func readLines(t *testing.T, path string) []string {
	t.Helper()
	dat, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var lines []string
	for _, l := range strings.Split(string(dat), "\n") {
		if l != "" {
			lines = append(lines, l)
		}
	}
	return lines
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	previous := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w
	defer func() {
		os.Stdout = previous
		_ = r.Close()
		_ = w.Close()
	}()

	fn()
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	output, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	return string(output)
}

// assertSet checks that got contains exactly the strings in want (order independent).
func assertSet(t *testing.T, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("got %d lines %v, want %d %v", len(got), got, len(want), want)
	}
	gotSet := make(map[string]struct{}, len(got))
	for _, g := range got {
		gotSet[g] = struct{}{}
	}
	wantSet := make(map[string]struct{}, len(want))
	for _, w := range want {
		wantSet[w] = struct{}{}
	}
	if len(gotSet) != len(got) {
		t.Errorf("got duplicate lines in %v", got)
	}
	if !reflect.DeepEqual(gotSet, wantSet) {
		t.Errorf("got set %v, want %v", gotSet, wantSet)
	}
}

func TestRunEmails(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join(dir, "in.txt")
	if err := os.WriteFile(in, []byte("a@x.com\nB@X.COM\na@x.com\n"), 0644); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(dir, "out.txt")
	if code := run([]string{"-input", in, "-output", out, "-count", "emails"}); code != 0 {
		t.Fatalf("run exited with %d", code)
	}
	assertSet(t, readLines(t, out), []string{"a@x.com", "b@x.com"})
}

func TestRunSha256Input(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join(dir, "in.txt")
	hash := "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"
	if err := os.WriteFile(in, []byte(hash+"\n"+hash+"\n"), 0644); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(dir, "out.txt")
	if code := run([]string{"-input", in, "-output", out, "-count", "sha256"}); code != 0 {
		t.Fatalf("run exited with %d", code)
	}
	assertSet(t, readLines(t, out), []string{hash})
}

func TestRunURLs(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join(dir, "in.txt")
	if err := os.WriteFile(in, []byte("https://example.com/a\nhttp://test.io/b\n"), 0644); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(dir, "out.txt")
	if code := run([]string{"-input", in, "-output", out, "-count", "urls"}); code != 0 {
		t.Fatalf("run exited with %d", code)
	}
	assertSet(t, readLines(t, out), []string{"https://example.com/a", "http://test.io/b"})
}

func TestRunDNI(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join(dir, "in.txt")
	if err := os.WriteFile(in, []byte("dni 12345678Z and 87654321x\n"), 0644); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(dir, "out.txt")
	if code := run([]string{"-input", in, "-output", out, "-count", "dnis"}); code != 0 {
		t.Fatalf("run exited with %d", code)
	}
	assertSet(t, readLines(t, out), []string{"12345678Z", "87654321X"})
}

func TestRunEncrypt(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join(dir, "in.txt")
	if err := os.WriteFile(in, []byte("A@X.COM\n"), 0644); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(dir, "out.txt")
	if code := run([]string{"-input", in, "-output", out, "-count", "emails", "-encrypt=true"}); code != 0 {
		t.Fatalf("run exited with %d", code)
	}
	hash := stringToSha256("a@x.com")
	assertSet(t, readLines(t, out), []string{hash})
}

func TestRunOutputInsideFolder(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join(dir, "in.txt")
	if err := os.WriteFile(in, []byte("a@x.com\n"), 0644); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(dir, "out.txt")
	if code := run([]string{"-input", dir, "-output", out, "-count", "emails"}); code != 1 {
		t.Fatalf("run exited with %d, want 1 (output inside folder rejected)", code)
	}
}

func TestRunHelp(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		wantHelpText bool
	}{
		{"short", []string{"-h"}, false},
		{"long", []string{"-help"}, true},
		{"double dash", []string{"--help"}, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var code int
			output := captureStdout(t, func() {
				code = run(tc.args)
			})
			if code != 0 {
				t.Fatalf("run %v exited with %d, want 0", tc.args, code)
			}
			if tc.wantHelpText && !strings.Contains(output, "* * * HELP * * *") {
				t.Errorf("run %v output %q, want help text", tc.args, output)
			}
		})
	}
}

func TestRunMissingInput(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") == "1" {
		run([]string{"-input", filepath.Join(t.TempDir(), "missing.txt"), "-output", filepath.Join(t.TempDir(), "out.txt")})
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestRunMissingInput")
	cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected non-zero exit, got output: %s", out)
	}
	if !strings.Contains(string(out), "ERROR: The file/path") {
		t.Fatalf("expected error message, got: %s", out)
	}
}
