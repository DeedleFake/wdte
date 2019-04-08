(window.webpackJsonp=window.webpackJsonp||[]).push([[0],{225:function(e,t,n){e.exports=n(479)},479:function(e,t,n){"use strict";n.r(t);var o={};n.r(o),n.d(o,"fib",function(){return C}),n.d(o,"stream",function(){return D}),n.d(o,"strings",function(){return F}),n.d(o,"lambdas",function(){return S}),n.d(o,"quine",function(){return q}),n.d(o,"hundredDoors",function(){return I});var i=n(1),a=n.n(i),r=n(30),s=n.n(r),c=(n(230),n(52)),l=n.n(c),u=n(80),d=n(53),m=n(120),p=n(25),h=n(124),f=n(192),g=n.n(f),b=n(193),w=n.n(b),v=n(78),x=n.n(v);n(319);x.a.define("ace/mode/wdte",function(e,t,n){var o=e("ace/lib/oop"),i=e("ace/mode/text").Mode,a=e("ace/mode/text_highlight_rules").TextHighlightRules,r=function(){this.$rules={start:[{token:"comment",regex:"^#.*$"},{token:"constant",regex:"[0-9]+|[0-9]+\\.[0-9]+|\\.[0-9]+"},{token:"keyword",regex:"(\\b(memo|let|import)\\b)"},{token:"keyword.operator",regex:"\\.|\\{|\\}|\\[|\\]|\\(|\\)|=>|;|:|->|--|-\\||\\(@"},{token:"string",regex:"(([\"']).*\\2)"}]}};o.inherits(r,a);var s=function(){this.HighlightRules=r};o.inherits(s,i),t.Mode=s});var y=n(194),k=n(487),E=n(488),O=n(486),j="Introduction\n============\n\nWelcome to the WDTE playground, a browser based evaluation environment for WDTE. This playground's features includes the standard function set as well as a number of importable modules.\n\nIf you have never seen WDTE before and are completely confused at the moment, try reading the overview on the WDTE project's wiki: https://github.com/DeedleFake/wdte/wiki\n\nFun Fact\n--------\n\nThe WDTE interpreter has been compiled to WebAssembly for this playground, meaning that, by opening this page, you've downloaded the entire system. Congratulations.\n\nDocumentation\n-------------\n\nFor documentation on the standard function set, see https://godoc.org/github.com/DeedleFake/wdte/std\n\nThe standard library is available for importing, with the exception of the `io/file` module. The `io` module is pre-inserted into the initial scope as `io`.",T="io.stdout -> io.writeln 'Greetings, pocket universe.';",C={name:"Fibonacci",desc:"\nFibonacci\n=========\n\nThis example provides a memoized implementation of a recursive Fibonacci number generator. It also provides a recursive factorial function for the heck of it.\n",input:"let memo fib n => n {>= 2 => + (fib (- n 1)) (fib (- n 2))};\n\nlet ! n => n {\n\t<= 1 => 1;\n\ttrue => - n 1 -> ! -> * n;\n};\n\nfib 30\n-- io.writeln io.stdout\n-> / 5\n-- io.writeln io.stdout\n;"},D={name:"Stream",desc:"\nStream\n======\n\nThis example demonstrates the `stream` module. This module provides functional iterator operations, such as map, reduce, and filter.\n\nFor a full list of functions, see [the godocs][godoc].\n\n[godoc]: https://www.godoc.org/github.com/DeedleFake/wdte/std/stream\n",input:"let m => import 'math';\nlet s => import 'stream';\n\nio.writeln io.stdout 'Map and filter:';\ns.range 0 (* m.pi 2) (/ m.pi 2)\n-> s.map m.sin\n-> s.filter (>= 0)\n-> s.map (io.writeln io.stdout)\n-> s.drain\n;\n\nio.writeln io.stdout 'Reduce:';\ns.range 1 5\n-> s.reduce 1 *\n-- io.writeln io.stdout\n;"},F={name:"Strings",desc:"\nStrings\n=======\n\nThis example demonstrates the `strings` module. This module provides basic string operations, such as finding the index of a substring, as well as more complicated operations, such as formatting.\n\nFor a full list of functions, including an explanation of the formatting system, see [the godocs][godoc].\n\n[godoc]: https://www.godoc.org/github.com/DeedleFake/wdte/std/strings\n",input:"let a => import 'arrays';\nlet s => import 'stream';\nlet str => import 'strings';\n\na.stream ['abc'; 'bcd'; 'cde']\n-> s.map (str.index 'cd')\n-> s.collect\n-- io.writeln io.stdout\n;\n\n'This is the type of English up with which I will not put.'\n-> str.format '{q}'\n-- io.writeln io.stdout\n;"},S={name:"Lambdas",desc:"\nLambdas\n=======\n\nThis example demonstrates lambdas by implementing an iterative Fibonacci number calculator using streams.\n",input:"let s => import 'stream';\nlet a => import 'arrays';\n\nlet fib n => s.range 1 n\n\t-> s.reduce [0; 1] (@ self p n =>\n\t\tlet [a b] => p;\n\t\t[\n\t\t\tb;\n\t\t\t+ a b;\n\t\t];\n\t)\n\t-> at 1\n\t;\n\nfib 30\n-- io.writeln io.stdout\n;"},q={name:"Quine",desc:"\nQuine\n=====\n\nThis example is an implemenation of a quine. That's about it.\n",input:"let str => import 'strings';\nlet q => \"let str => import 'strings';\\nlet q => {q};\\nstr.format q q -- io.writeln io.stdout;\";\nstr.format q q -- io.writeln io.stdout;"},I={name:"100 Doors",desc:"\n100 Doors\n=========\n\nThe [100 doors problem](https://www.rosettacode.org/wiki/100_doors), as presented by Rosetta Code, is as follows:\n\nThere are 100 doors that are all closed. You walk past the doors in the same direction 100 times. On the first pass, you toggle the state of every door, opening closed doors and closing open doors. On the second pass, you toggle every second door. On the third you toggle every third door. Etc.\n\nThis example simulates this scenario, printing out the final state of the doors.\n",input:"let a => import 'arrays';\nlet s => import 'stream';\n\nlet toggle doors m =>\n\ta.stream doors\n\t-> s.enumerate\n\t-> s.map (@ s n => [+ (at n 0) 1; at n 1])\n\t-> s.map (@ s n => n {\n\t\t\t(@ s n => == (% (at n 0) m) 0) => ! (at n 1);\n\t\t\ttrue => at n 1;\n\t\t})\n\t-> s.collect\n\t;\n\ns.range 100\n-> s.map false\n-> s.collect : doors\n-> s.range 1 100\n-> s.reduce doors toggle\n-> a.stream\n-> s.map (@ s n => 0 {\n\t\tn => 'Open';\n\t\ttrue => 'Closed';\n\t} -- io.writeln io.stdout)\n-> s.drain\n;"},R=function(){var e=Object(u.a)(l.a.mark(function e(t){return l.a.wrap(function(e){for(;;)switch(e.prev=e.next){case 0:return e.next=2,new Promise(function(e,n){window.WDTE.run(t,function(t,o){if(null!=t)return n(t);e(o)})});case 2:return e.abrupt("return",e.sent);case 3:case"end":return e.stop()}},e,this)}));return function(t){return e.apply(this,arguments)}}(),A=Object(y.a)(function(e){return{"@font-face":{fontFamily:"Go Mono",src:"url(assets/Go-Mono.ttf)"},main:{display:"flex",flexDirection:"row",backgroundColor:"#EEEEEE",boxSizing:"border-box",padding:8,width:"100%",height:"100%",position:"absolute"},column:{display:"flex",flexDirection:"column",flex:1,margin:8,overflowY:"auto",minWidth:300},message:{marginBottom:"8px !important"},output:{minHeight:300,fontFamily:"Go-Mono",fontSize:12,flex:1,overflow:"auto",padding:8,boxShadow:"inset 4px 4px 4px #AAAAAA",borderRadius:8,backgroundColor:"#CCCCCC"},slide:{"&.enter":{transition:"all 300ms",overflow:"hidden",maxHeight:0,"&.active":{maxHeight:"500px"}},"&.exit":{transition:"all 300ms",overflow:"hidden",maxHeight:"500px","&.active":{maxHeight:0}}}}}),W=function(e){var t=A(),n=Object(i.useState)(j),r=Object(p.a)(n,2),s=r[0],c=r[1],f=Object(i.useState)(function(){return""!==window.location.hash?decodeURIComponent(window.location.hash.substr(1)):T}),b=Object(p.a)(f,2),v=b[0],x=b[1],y=Object(i.useState)(""),C=Object(p.a)(y,2),D=C[0],F=C[1],S=Object(i.useMemo)(function(){return encodeURIComponent(v)},[v]),q=Object(i.useReducer)(function(e,t){switch(t.$){case"add":return null!=t.timeout&&setTimeout(function(){H({$:"remove",id:e.id})},t.timeout),Object(m.a)({},e,Object(d.a)({id:e.id+1},e.id,{type:t.type,msg:t.msg}));case"remove":return Object.entries(e).filter(function(e){var n=Object(p.a)(e,2),o=n[0];return n[1],o!==t.id.toString()}).reduce(function(e,t){var n=Object(p.a)(t,2),o=n[0],i=n[1];return Object(m.a)({},e,Object(d.a)({},o,i))},{});default:return e}},{id:0}),I=Object(p.a)(q,2),W=I[0],H=I[1],M=Object(i.useCallback)(function(e,t){var n=arguments.length>2&&void 0!==arguments[2]?arguments[2]:3e3;H({$:"add",type:e,msg:t,timeout:n})},[W]),N=Object(i.useCallback)(Object(u.a)(l.a.mark(function e(){return l.a.wrap(function(e){for(;;)switch(e.prev=e.next){case 0:return e.prev=0,e.t0=F,e.next=4,R(v);case 4:e.t1=e.sent,(0,e.t0)(e.t1),e.next=11;break;case 8:e.prev=8,e.t2=e.catch(0),F(e.t2.toString());case 11:case"end":return e.stop()}},e,this,[[0,8]])})),[v]),G=Object(i.useCallback)(function(){try{!function(e){var t=document.createElement("textarea");t.value=e;try{if(document.body.appendChild(t),t.focus(),t.select(),!document.execCommand("copy"))throw new Error("copy failed")}finally{document.body.removeChild(t)}}("".concat(window.location.origin).concat(window.location.pathname,"#").concat(S)),M("success","Link successfully copied to clipboard.")}catch(e){M("error","Failed to copy to clipboard: ".concat(e.toString()))}},[S]);return a.a.createElement("div",{className:t.main},a.a.createElement(h.TransitionGroup,{component:"div",className:t.column},Object.entries(W).filter(function(e){var t=Object(p.a)(e,2),n=t[0];return t[1],!isNaN(n)}).map(function(e){var n=Object(p.a)(e,2),o=n[0],i=n[1];return a.a.createElement(h.CSSTransition,{key:o,classNames:{enter:"enter",enterActive:"active",exit:"exit",exitActive:"active"},timeout:300},a.a.createElement(k.a,Object.assign({className:[t.message,t.slide].join(" ")},Object(d.a)({},i.type,!0)),a.a.createElement("p",null,i.msg)))}),a.a.createElement(g.a,{source:s})),a.a.createElement("div",{className:t.column},a.a.createElement(E.a,{inverted:!0},a.a.createElement(E.a.Item,{onClick:N},"Run"),a.a.createElement(O.a,{item:!0,text:"Examples"},a.a.createElement(O.a.Menu,null,Object.entries(o).map(function(e){var t=Object(p.a)(e,2),n=t[0],i=t[1];return a.a.createElement(O.a.Item,{key:n,value:n,onClick:function(e,t){c(o[t.value].desc),x(o[t.value].input)}},i.name)}))),a.a.createElement(E.a.Item,{position:"right",onClick:G},"Share")),a.a.createElement(w.a,{style:{width:null,height:null,minHeight:300,flex:1,borderRadius:8},mode:"wdte",theme:"vibrant_ink",value:v,onChange:function(e){return x(e)}}),a.a.createElement("pre",{className:t.output},D)))};s.a.render(a.a.createElement(W,null),document.getElementById("root"))}},[[225,1,2]]]);
//# sourceMappingURL=main.71037acc.chunk.js.map