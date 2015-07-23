package suffuse

import (
  "strings"
  "fmt"
  . "gopkg.in/check.v1"
)

func (s *Tsfs) TestStrings(c *C) {
  text := "foo\nbar\n"
  s1 := Exec("echo", strings.TrimSpace(text))
  c.Assert(s1.Lines(), DeepEquals, Strings{"foo", "bar"})
  c.Assert(s1.Slurp(), Equals, text)
}

func (s *Tsfs) TestStripMargin(c *C) {
  s1 := StripMargin('|', `
    |the quick
    |brown fox
  `)
  c.Assert(s1, Equals, "the quick\nbrown fox")
}

func (s *Tsfs) TestStringMethods(c *C) {
    ls := Strings{"dog  ", "  cat", " monkey in middle "}
   one := ls.JoinWords()
    fm := ls.FlatMap(func(s string)[]string { return []string{"a", "b"} }).JoinWords()
  fold := ls.Fold(func(acc, s string)string { return fmt.Sprintf("%s%d!", acc, len(s)) })

  c.Assert(one, Equals, "dog cat monkey in middle")
  c.Assert(fm, Equals, "a b a b a b")
  c.Assert(fold, Equals, "5!5!18!")
}
