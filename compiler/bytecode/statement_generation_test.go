package bytecode

import "testing"

func TestWhileStatementInBlock(t *testing.T) {
	input := `
	i = 1
	thread do
	  puts(i)
	  while i <= 1000 do
		puts(i)
		i = i + 1
	  end
	end
	`

	expected := `
<Block:0>
0 putself
1 getlocal 1 0
2 send puts 1
3 jump 14
4 putnil
5 pop
6 jump 14
7 putself
8 getlocal 1 0
9 send puts 1
10 getlocal 1 0
11 putobject 1
12 send + 1
13 setlocal 1 0
14 getlocal 1 0
15 putobject 1000
16 send <= 1
17 branchif 7
18 putnil
19 pop
20 leave
<ProgramStart>
0 putobject 1
1 setlocal 0 0
2 putself
3 send thread 0 block:0
4 leave
`

	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}

func TestNextStatement(t *testing.T) {
	input := `
	x = 0
	y = 0

	while x < 10 do
	  x = x + 1
	  if x == 5
	    next
	  end
	  y = y + 1
	end
	`

	expected := `
<ProgramStart>
0 putobject 0
1 setlocal 0 0
2 putobject 0
3 setlocal 0 1
4 jump 21
5 putnil
6 pop
7 jump 21
8 getlocal 0 0
9 putobject 1
10 send + 1
11 setlocal 0 0
12 getlocal 0 0
13 putobject 5
14 send == 1
15 branchunless 17
16 jump 21
17 getlocal 0 1
18 putobject 1
19 send + 1
20 setlocal 0 1
21 getlocal 0 0
22 putobject 10
23 send < 1
24 branchif 8
25 putnil
26 pop
27 leave
`

	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}

func TestNamespacedClass(t *testing.T) {
	input := `
	module Foo
	  class Bar
	    class Baz
	      def bar
	      end
	    end
	  end
	end

	Foo::Bar::Baz.new.bar
	`

	expected := `
<Def:Object::Foo::Bar::Baz#bar>
0 putnil
1 leave
<DefClass:Baz>
0 putself
1 putstring "Object::Foo::Bar::Baz#bar"
2 def_method 0
3 leave
<DefClass:Bar>
0 putself
1 def_class class:Baz
2 pop
3 leave
<DefClass:Foo>
0 putself
1 def_class class:Bar
2 pop
3 leave
<ProgramStart>
0 putself
1 def_class module:Foo
2 pop
3 getconstant Foo
4 getconstant Bar
5 getconstant Baz
6 send new 0
7 send bar 0
8 leave
`

	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}

func TestClassMethodDefinition(t *testing.T) {
	input := `
class Foo
  def self.bar
    10
  end
end

Foo.bar
`
	expected := `
<Def:Object::Foo.bar>
0 putobject 10
1 leave
<DefClass:Foo>
0 putself
1 putstring "Object::Foo.bar"
2 def_singleton_method 0
3 leave
<ProgramStart>
0 putself
1 def_class class:Foo
2 pop
3 getconstant Foo
4 send bar 0
5 leave
`

	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}

func TestClassCompilation(t *testing.T) {
	input := `
class Bar
  def bar
    10
  end
end

class Foo < Bar
end

Foo.new.bar
`
	expected := `
<Def:Object::Bar#bar>
0 putobject 10
1 leave
<DefClass:Bar>
0 putself
1 putstring "Object::Bar#bar"
2 def_method 0
3 leave
<DefClass:Foo>
0 leave
<ProgramStart>
0 putself
1 def_class class:Bar
2 pop
3 putself
4 getconstant Bar
5 def_class class:Foo Bar
6 pop
7 getconstant Foo
8 send new 0
9 send bar 0
10 leave
`
	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}

func TestModuleCompilation(t *testing.T) {
	input := `


module Bar
  def bar
    10
  end
end

class Foo
  include(Bar)
end

Foo.new.bar
`
	expected := `
<Def:Object::Bar#bar>
0 putobject 10
1 leave
<DefClass:Bar>
0 putself
1 putstring "Object::Bar#bar"
2 def_method 0
3 leave
<DefClass:Foo>
0 putself
1 getconstant Bar
2 send include 1
3 leave
<ProgramStart>
0 putself
1 def_class module:Bar
2 pop
3 putself
4 def_class class:Foo
5 pop
6 getconstant Foo
7 send new 0
8 send bar 0
9 leave
`
	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}

func TestWhileStatementWithoutMethodCallInCondition(t *testing.T) {
	input := `
	i = 10

	while i > 0 do
	  i = i - 1
	end

	i
`
	expected := `
<ProgramStart>
0 putobject 10
1 setlocal 0 0
2 jump 10
3 putnil
4 pop
5 jump 10
6 getlocal 0 0
7 putobject 1
8 send - 1
9 setlocal 0 0
10 getlocal 0 0
11 putobject 0
12 send > 1
13 branchif 6
14 putnil
15 pop
16 getlocal 0 0
17 leave
`

	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}

func TestWhileStatementWithMethodCallInCondition(t *testing.T) {
	input := `
	i = 10
	a = [1, 2, 3]

	while i > a.length do
	  i = i - 1
	end

	i
`
	expected := `
<ProgramStart>
0 putobject 10
1 setlocal 0 0
2 putobject 1
3 putobject 2
4 putobject 3
5 newarray 3
6 setlocal 0 1
7 jump 15
8 putnil
9 pop
10 jump 15
11 getlocal 0 0
12 putobject 1
13 send - 1
14 setlocal 0 0
15 getlocal 0 0
16 getlocal 0 1
17 send length 0
18 send > 1
19 branchif 11
20 putnil
21 pop
22 getlocal 0 0
23 leave
`

	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}

func TestRemoveUnusedExpression(t *testing.T) {
	input := `
	i = 0

	while i < 100 do
	  10
	  i++
	end

	i
	`

	expected := `
<ProgramStart>
0 putobject 0
1 setlocal 0 0
2 jump 9
3 putnil
4 pop
5 jump 9
6 getlocal 0 0
7 send ++ 0
8 pop
9 getlocal 0 0
10 putobject 100
11 send < 1
12 branchif 6
13 putnil
14 pop
15 getlocal 0 0
16 leave
`

	bytecode := compileToBytecode(input)
	compareBytecode(t, bytecode, expected)
}
