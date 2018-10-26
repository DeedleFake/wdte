// Code generated by pgen. DO NOT EDIT.

package pgen

var Table = map[Lookup]Rule{
	{Term: newTerm(";"), NTerm: newNTerm("aexprs")}:        newRule(newTerm(";")),
	{Term: newTerm("number"), NTerm: newNTerm("aexprs")}:   newRule(newNTerm("exprs")),
	{Term: newTerm("string"), NTerm: newNTerm("aexprs")}:   newRule(newNTerm("exprs")),
	{Term: newTerm("["), NTerm: newNTerm("aexprs")}:        newRule(newNTerm("exprs")),
	{Term: newTerm("id"), NTerm: newNTerm("aexprs")}:       newRule(newNTerm("exprs")),
	{Term: newTerm("("), NTerm: newNTerm("aexprs")}:        newRule(newNTerm("exprs")),
	{Term: newTerm("import"), NTerm: newNTerm("aexprs")}:   newRule(newNTerm("exprs")),
	{Term: newTerm("(@"), NTerm: newNTerm("aexprs")}:       newRule(newNTerm("exprs")),
	{Term: newTerm("]"), NTerm: newNTerm("aexprs")}:        newRule(newEpsilon()),
	{Term: newTerm("id"), NTerm: newNTerm("argdecls")}:     newRule(newTerm("id"), newNTerm("argdecls")),
	{Term: newTerm("=>"), NTerm: newNTerm("argdecls")}:     newRule(newEpsilon()),
	{Term: newTerm("("), NTerm: newNTerm("args")}:          newRule(newNTerm("single"), newNTerm("args")),
	{Term: newTerm("number"), NTerm: newNTerm("args")}:     newRule(newNTerm("single"), newNTerm("args")),
	{Term: newTerm("string"), NTerm: newNTerm("args")}:     newRule(newNTerm("single"), newNTerm("args")),
	{Term: newTerm("import"), NTerm: newNTerm("args")}:     newRule(newNTerm("single"), newNTerm("args")),
	{Term: newTerm("(@"), NTerm: newNTerm("args")}:         newRule(newNTerm("single"), newNTerm("args")),
	{Term: newTerm("["), NTerm: newNTerm("args")}:          newRule(newNTerm("single"), newNTerm("args")),
	{Term: newTerm("id"), NTerm: newNTerm("args")}:         newRule(newNTerm("single"), newNTerm("args")),
	{Term: newTerm(";"), NTerm: newNTerm("args")}:          newRule(newEpsilon()),
	{Term: newTerm("=>"), NTerm: newNTerm("args")}:         newRule(newEpsilon()),
	{Term: newTerm(":"), NTerm: newNTerm("args")}:          newRule(newEpsilon()),
	{Term: newTerm("{"), NTerm: newNTerm("args")}:          newRule(newEpsilon()),
	{Term: newTerm("->"), NTerm: newNTerm("args")}:         newRule(newEpsilon()),
	{Term: newTerm("--"), NTerm: newNTerm("args")}:         newRule(newEpsilon()),
	{Term: newTerm("["), NTerm: newNTerm("array")}:         newRule(newTerm("["), newNTerm("aexprs"), newTerm("]")),
	{Term: newTerm("("), NTerm: newNTerm("cexprs")}:        newRule(newNTerm("expr"), newTerm(";"), newNTerm("cexprs")),
	{Term: newTerm("id"), NTerm: newNTerm("cexprs")}:       newRule(newNTerm("expr"), newTerm(";"), newNTerm("cexprs")),
	{Term: newTerm("let"), NTerm: newNTerm("cexprs")}:      newRule(newNTerm("letexpr"), newTerm(";"), newNTerm("cexprs")),
	{Term: newTerm("number"), NTerm: newNTerm("cexprs")}:   newRule(newNTerm("expr"), newTerm(";"), newNTerm("cexprs")),
	{Term: newTerm("string"), NTerm: newNTerm("cexprs")}:   newRule(newNTerm("expr"), newTerm(";"), newNTerm("cexprs")),
	{Term: newTerm("import"), NTerm: newNTerm("cexprs")}:   newRule(newNTerm("expr"), newTerm(";"), newNTerm("cexprs")),
	{Term: newTerm("(@"), NTerm: newNTerm("cexprs")}:       newRule(newNTerm("expr"), newTerm(";"), newNTerm("cexprs")),
	{Term: newTerm("["), NTerm: newNTerm("cexprs")}:        newRule(newNTerm("expr"), newTerm(";"), newNTerm("cexprs")),
	{Term: newTerm(")"), NTerm: newNTerm("cexprs")}:        newRule(newEpsilon()),
	{Term: newEOF(), NTerm: newNTerm("cexprs")}:            newRule(newEpsilon()),
	{Term: newTerm("--"), NTerm: newNTerm("chain")}:        newRule(newTerm("--"), newNTerm("expr")),
	{Term: newTerm("->"), NTerm: newNTerm("chain")}:        newRule(newTerm("->"), newNTerm("expr")),
	{Term: newTerm(";"), NTerm: newNTerm("chain")}:         newRule(newEpsilon()),
	{Term: newTerm("=>"), NTerm: newNTerm("chain")}:        newRule(newEpsilon()),
	{Term: newTerm("("), NTerm: newNTerm("compound")}:      newRule(newTerm("("), newNTerm("cexprs"), newTerm(")")),
	{Term: newTerm("["), NTerm: newNTerm("expr")}:          newRule(newNTerm("single"), newNTerm("args"), newNTerm("slot"), newNTerm("switch"), newNTerm("chain")),
	{Term: newTerm("id"), NTerm: newNTerm("expr")}:         newRule(newNTerm("single"), newNTerm("args"), newNTerm("slot"), newNTerm("switch"), newNTerm("chain")),
	{Term: newTerm("("), NTerm: newNTerm("expr")}:          newRule(newNTerm("single"), newNTerm("args"), newNTerm("slot"), newNTerm("switch"), newNTerm("chain")),
	{Term: newTerm("number"), NTerm: newNTerm("expr")}:     newRule(newNTerm("single"), newNTerm("args"), newNTerm("slot"), newNTerm("switch"), newNTerm("chain")),
	{Term: newTerm("string"), NTerm: newNTerm("expr")}:     newRule(newNTerm("single"), newNTerm("args"), newNTerm("slot"), newNTerm("switch"), newNTerm("chain")),
	{Term: newTerm("import"), NTerm: newNTerm("expr")}:     newRule(newNTerm("single"), newNTerm("args"), newNTerm("slot"), newNTerm("switch"), newNTerm("chain")),
	{Term: newTerm("(@"), NTerm: newNTerm("expr")}:         newRule(newNTerm("single"), newNTerm("args"), newNTerm("slot"), newNTerm("switch"), newNTerm("chain")),
	{Term: newTerm("id"), NTerm: newNTerm("exprs")}:        newRule(newNTerm("expr"), newTerm(";"), newNTerm("exprs")),
	{Term: newTerm("("), NTerm: newNTerm("exprs")}:         newRule(newNTerm("expr"), newTerm(";"), newNTerm("exprs")),
	{Term: newTerm("number"), NTerm: newNTerm("exprs")}:    newRule(newNTerm("expr"), newTerm(";"), newNTerm("exprs")),
	{Term: newTerm("string"), NTerm: newNTerm("exprs")}:    newRule(newNTerm("expr"), newTerm(";"), newNTerm("exprs")),
	{Term: newTerm("import"), NTerm: newNTerm("exprs")}:    newRule(newNTerm("expr"), newTerm(";"), newNTerm("exprs")),
	{Term: newTerm("(@"), NTerm: newNTerm("exprs")}:        newRule(newNTerm("expr"), newTerm(";"), newNTerm("exprs")),
	{Term: newTerm("["), NTerm: newNTerm("exprs")}:         newRule(newNTerm("expr"), newTerm(";"), newNTerm("exprs")),
	{Term: newTerm("]"), NTerm: newNTerm("exprs")}:         newRule(newEpsilon()),
	{Term: newTerm("memo"), NTerm: newNTerm("funcmod")}:    newRule(newTerm("memo")),
	{Term: newTerm("memo"), NTerm: newNTerm("funcmods")}:   newRule(newNTerm("funcmod"), newNTerm("funcmods")),
	{Term: newTerm("id"), NTerm: newNTerm("funcmods")}:     newRule(newEpsilon()),
	{Term: newTerm("import"), NTerm: newNTerm("import")}:   newRule(newTerm("import"), newTerm("string")),
	{Term: newTerm("(@"), NTerm: newNTerm("lambda")}:       newRule(newTerm("(@"), newNTerm("funcmods"), newTerm("id"), newNTerm("argdecls"), newTerm("=>"), newNTerm("cexprs"), newTerm(")")),
	{Term: newTerm("let"), NTerm: newNTerm("letexpr")}:     newRule(newTerm("let"), newNTerm("funcmods"), newTerm("id"), newNTerm("argdecls"), newTerm("=>"), newNTerm("expr")),
	{Term: newTerm("id"), NTerm: newNTerm("script")}:       newRule(newNTerm("cexprs"), newEOF()),
	{Term: newTerm("("), NTerm: newNTerm("script")}:        newRule(newNTerm("cexprs"), newEOF()),
	{Term: newEOF(), NTerm: newNTerm("script")}:            newRule(newNTerm("cexprs"), newEOF()),
	{Term: newTerm("let"), NTerm: newNTerm("script")}:      newRule(newNTerm("cexprs"), newEOF()),
	{Term: newTerm("string"), NTerm: newNTerm("script")}:   newRule(newNTerm("cexprs"), newEOF()),
	{Term: newTerm("["), NTerm: newNTerm("script")}:        newRule(newNTerm("cexprs"), newEOF()),
	{Term: newTerm("number"), NTerm: newNTerm("script")}:   newRule(newNTerm("cexprs"), newEOF()),
	{Term: newTerm("(@"), NTerm: newNTerm("script")}:       newRule(newNTerm("cexprs"), newEOF()),
	{Term: newTerm("import"), NTerm: newNTerm("script")}:   newRule(newNTerm("cexprs"), newEOF()),
	{Term: newTerm("import"), NTerm: newNTerm("single")}:   newRule(newNTerm("import")),
	{Term: newTerm("(@"), NTerm: newNTerm("single")}:       newRule(newNTerm("lambda")),
	{Term: newTerm("["), NTerm: newNTerm("single")}:        newRule(newNTerm("array")),
	{Term: newTerm("("), NTerm: newNTerm("single")}:        newRule(newNTerm("subbable")),
	{Term: newTerm("id"), NTerm: newNTerm("single")}:       newRule(newNTerm("subbable")),
	{Term: newTerm("number"), NTerm: newNTerm("single")}:   newRule(newTerm("number")),
	{Term: newTerm("string"), NTerm: newNTerm("single")}:   newRule(newTerm("string")),
	{Term: newTerm(":"), NTerm: newNTerm("slot")}:          newRule(newTerm(":"), newTerm("id")),
	{Term: newTerm("{"), NTerm: newNTerm("slot")}:          newRule(newEpsilon()),
	{Term: newTerm("->"), NTerm: newNTerm("slot")}:         newRule(newEpsilon()),
	{Term: newTerm("--"), NTerm: newNTerm("slot")}:         newRule(newEpsilon()),
	{Term: newTerm(";"), NTerm: newNTerm("slot")}:          newRule(newEpsilon()),
	{Term: newTerm("=>"), NTerm: newNTerm("slot")}:         newRule(newEpsilon()),
	{Term: newTerm("."), NTerm: newNTerm("sub")}:           newRule(newTerm("."), newNTerm("subbable")),
	{Term: newTerm("string"), NTerm: newNTerm("sub")}:      newRule(newEpsilon()),
	{Term: newTerm("id"), NTerm: newNTerm("sub")}:          newRule(newEpsilon()),
	{Term: newTerm("=>"), NTerm: newNTerm("sub")}:          newRule(newEpsilon()),
	{Term: newTerm("number"), NTerm: newNTerm("sub")}:      newRule(newEpsilon()),
	{Term: newTerm(":"), NTerm: newNTerm("sub")}:           newRule(newEpsilon()),
	{Term: newTerm("(@"), NTerm: newNTerm("sub")}:          newRule(newEpsilon()),
	{Term: newTerm("import"), NTerm: newNTerm("sub")}:      newRule(newEpsilon()),
	{Term: newTerm("("), NTerm: newNTerm("sub")}:           newRule(newEpsilon()),
	{Term: newTerm("->"), NTerm: newNTerm("sub")}:          newRule(newEpsilon()),
	{Term: newTerm("--"), NTerm: newNTerm("sub")}:          newRule(newEpsilon()),
	{Term: newTerm("{"), NTerm: newNTerm("sub")}:           newRule(newEpsilon()),
	{Term: newTerm("["), NTerm: newNTerm("sub")}:           newRule(newEpsilon()),
	{Term: newTerm(";"), NTerm: newNTerm("sub")}:           newRule(newEpsilon()),
	{Term: newTerm("id"), NTerm: newNTerm("subbable")}:     newRule(newTerm("id"), newNTerm("sub")),
	{Term: newTerm("("), NTerm: newNTerm("subbable")}:      newRule(newNTerm("compound"), newNTerm("sub")),
	{Term: newTerm("{"), NTerm: newNTerm("switch")}:        newRule(newTerm("{"), newNTerm("switches"), newTerm("}")),
	{Term: newTerm("->"), NTerm: newNTerm("switch")}:       newRule(newEpsilon()),
	{Term: newTerm("--"), NTerm: newNTerm("switch")}:       newRule(newEpsilon()),
	{Term: newTerm(";"), NTerm: newNTerm("switch")}:        newRule(newEpsilon()),
	{Term: newTerm("=>"), NTerm: newNTerm("switch")}:       newRule(newEpsilon()),
	{Term: newTerm("id"), NTerm: newNTerm("switches")}:     newRule(newNTerm("expr"), newTerm("=>"), newNTerm("expr"), newTerm(";"), newNTerm("switches")),
	{Term: newTerm("("), NTerm: newNTerm("switches")}:      newRule(newNTerm("expr"), newTerm("=>"), newNTerm("expr"), newTerm(";"), newNTerm("switches")),
	{Term: newTerm("number"), NTerm: newNTerm("switches")}: newRule(newNTerm("expr"), newTerm("=>"), newNTerm("expr"), newTerm(";"), newNTerm("switches")),
	{Term: newTerm("string"), NTerm: newNTerm("switches")}: newRule(newNTerm("expr"), newTerm("=>"), newNTerm("expr"), newTerm(";"), newNTerm("switches")),
	{Term: newTerm("import"), NTerm: newNTerm("switches")}: newRule(newNTerm("expr"), newTerm("=>"), newNTerm("expr"), newTerm(";"), newNTerm("switches")),
	{Term: newTerm("(@"), NTerm: newNTerm("switches")}:     newRule(newNTerm("expr"), newTerm("=>"), newNTerm("expr"), newTerm(";"), newNTerm("switches")),
	{Term: newTerm("["), NTerm: newNTerm("switches")}:      newRule(newNTerm("expr"), newTerm("=>"), newNTerm("expr"), newTerm(";"), newNTerm("switches")),
	{Term: newTerm("}"), NTerm: newNTerm("switches")}:      newRule(newEpsilon()),
}
