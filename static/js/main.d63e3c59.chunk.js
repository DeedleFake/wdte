(window["webpackJsonpwdte-playground"]=window["webpackJsonpwdte-playground"]||[]).push([[0],{278:function(e,t,n){e.exports=n(603)},514:function(e,t,n){"use strict";(function(e,t){var o=n(40),i=n.n(o),r=n(66),a=n(265),s=n(266),c=n(29);!function(){if("undefined"!==typeof e);else{if("undefined"===typeof window)throw new Error("cannot export Go (neither global, window nor self is defined)");window.global=window}var o=e.process&&"node"===e.process.title;if(o){e.fs=n(515);var u=n(516);e.crypto={getRandomValues:function(e){u.randomFillSync(e)}},e.performance={now:function(){var e=t.hrtime(),n=Object(c.a)(e,2);return 1e3*n[0]+n[1]/1e6}};var l=n(599);e.TextEncoder=l.TextEncoder,e.TextDecoder=l.TextDecoder}else{var d="";e.fs={constants:{O_WRONLY:-1,O_RDWR:-1,O_CREAT:-1,O_TRUNC:-1,O_APPEND:-1,O_EXCL:-1},writeSync:function(e,t){var n=(d+=f.decode(t)).lastIndexOf("\n");return-1!==n&&(console.log(d.substr(0,n)),d=d.substr(n+1)),t.length},write:function(e,t,n,o,i,r){if(0!==n||o!==t.length||null!==i)throw new Error("not implemented");r(null,this.writeSync(e,t))},open:function(e,t,n,o){var i=new Error("not implemented");i.code="ENOSYS",o(i)},read:function(e,t,n,o,i,r){var a=new Error("not implemented");a.code="ENOSYS",r(a)},fsync:function(e,t){t(null)}}}var m=new TextEncoder("utf-8"),f=new TextDecoder("utf-8");if(e.Go=function(){function t(){var n=this;Object(a.a)(this,t),this.argv=["js"],this.env={},this.exit=function(e){0!==e&&console.warn("exit code:",e)},this._exitPromise=new Promise(function(e){n._resolveExitPromise=e}),this._pendingEvent=null,this._scheduledTimeouts=new Map,this._nextCallbackTimeoutID=1;var o=function(){return new DataView(n._inst.exports.mem.buffer)},i=function(e,t){o().setUint32(e+0,t,!0),o().setUint32(e+4,Math.floor(t/4294967296),!0)},r=function(e){return o().getUint32(e+0,!0)+4294967296*o().getInt32(e+4,!0)},s=function(e){var t=o().getFloat64(e,!0);if(0!==t){if(!isNaN(t))return t;var i=o().getUint32(e,!0);return n._values[i]}},c=function(e,t){if("number"===typeof t)return isNaN(t)?(o().setUint32(e+4,2146959360,!0),void o().setUint32(e,0,!0)):0===t?(o().setUint32(e+4,2146959360,!0),void o().setUint32(e,1,!0)):void o().setFloat64(e,t,!0);switch(t){case void 0:return void o().setFloat64(e,0,!0);case null:return o().setUint32(e+4,2146959360,!0),void o().setUint32(e,2,!0);case!0:return o().setUint32(e+4,2146959360,!0),void o().setUint32(e,3,!0);case!1:return o().setUint32(e+4,2146959360,!0),void o().setUint32(e,4,!0)}var i=n._refs.get(t);void 0===i&&(i=n._values.length,n._values.push(t),n._refs.set(t,i));var r=0;switch(typeof t){case"string":r=1;break;case"symbol":r=2;break;case"function":r=3}o().setUint32(e+4,2146959360|r,!0),o().setUint32(e,i,!0)},u=function(e){var t=r(e+0),o=r(e+8);return new Uint8Array(n._inst.exports.mem.buffer,t,o)},l=function(e){for(var t=r(e+0),n=r(e+8),o=new Array(n),i=0;i<n;i++)o[i]=s(t+8*i);return o},d=function(e){var t=r(e+0),o=r(e+8);return f.decode(new DataView(n._inst.exports.mem.buffer,t,o))},p=Date.now()-performance.now();this.importObject={go:{"runtime.wasmExit":function(e){var t=o().getInt32(e+8,!0);n.exited=!0,delete n._inst,delete n._values,delete n._refs,n.exit(t)},"runtime.wasmWrite":function(t){var i=r(t+8),a=r(t+16),s=o().getInt32(t+24,!0);e.fs.writeSync(i,new Uint8Array(n._inst.exports.mem.buffer,a,s))},"runtime.nanotime":function(e){i(e+8,1e6*(p+performance.now()))},"runtime.walltime":function(e){var t=(new Date).getTime();i(e+8,t/1e3),o().setInt32(e+16,t%1e3*1e6,!0)},"runtime.scheduleTimeoutEvent":function(e){var t=n._nextCallbackTimeoutID;n._nextCallbackTimeoutID++,n._scheduledTimeouts.set(t,setTimeout(function(){n._resume()},r(e+8)+1)),o().setInt32(e+16,t,!0)},"runtime.clearTimeoutEvent":function(e){var t=o().getInt32(e+8,!0);clearTimeout(n._scheduledTimeouts.get(t)),n._scheduledTimeouts.delete(t)},"runtime.getRandomData":function(e){crypto.getRandomValues(u(e+8))},"syscall/js.stringVal":function(e){c(e+24,d(e+8))},"syscall/js.valueGet":function(e){var t=Reflect.get(s(e+8),d(e+16));e=n._inst.exports.getsp(),c(e+32,t)},"syscall/js.valueSet":function(e){Reflect.set(s(e+8),d(e+16),s(e+32))},"syscall/js.valueIndex":function(e){c(e+24,Reflect.get(s(e+8),r(e+16)))},"syscall/js.valueSetIndex":function(e){Reflect.set(s(e+8),r(e+16),s(e+24))},"syscall/js.valueCall":function(e){try{var t=s(e+8),i=Reflect.get(t,d(e+16)),r=l(e+32),a=Reflect.apply(i,t,r);e=n._inst.exports.getsp(),c(e+56,a),o().setUint8(e+64,1)}catch(u){c(e+56,u),o().setUint8(e+64,0)}},"syscall/js.valueInvoke":function(e){try{var t=s(e+8),i=l(e+16),r=Reflect.apply(t,void 0,i);e=n._inst.exports.getsp(),c(e+40,r),o().setUint8(e+48,1)}catch(a){c(e+40,a),o().setUint8(e+48,0)}},"syscall/js.valueNew":function(e){try{var t=s(e+8),i=l(e+16),r=Reflect.construct(t,i);e=n._inst.exports.getsp(),c(e+40,r),o().setUint8(e+48,1)}catch(a){c(e+40,a),o().setUint8(e+48,0)}},"syscall/js.valueLength":function(e){i(e+16,parseInt(s(e+8).length))},"syscall/js.valuePrepareString":function(e){var t=m.encode(String(s(e+8)));c(e+16,t),i(e+24,t.length)},"syscall/js.valueLoadString":function(e){var t=s(e+8);u(e+16).set(t)},"syscall/js.valueInstanceOf":function(e){o().setUint8(e+24,s(e+8)instanceof s(e+16))},debug:function(e){console.log(e)}}}}return Object(s.a)(t,[{key:"run",value:function(){var t=Object(r.a)(i.a.mark(function t(n){var o,r,a,s,c,u,l,d=this;return i.a.wrap(function(t){for(;;)switch(t.prev=t.next){case 0:return this._inst=n,this._values=[NaN,0,null,!0,!1,e,this._inst.exports.mem,this],this._refs=new Map,this.exited=!1,o=new DataView(this._inst.exports.mem.buffer),r=4096,a=function(e){var t=r;return new Uint8Array(o.buffer,r,e.length+1).set(m.encode(e+"\0")),r+=e.length+(8-e.length%8),t},s=this.argv.length,c=[],this.argv.forEach(function(e){c.push(a(e))}),u=Object.keys(this.env).sort(),c.push(u.length),u.forEach(function(e){c.push(a("".concat(e,"=").concat(d.env[e])))}),l=r,c.forEach(function(e){o.setUint32(r,e,!0),o.setUint32(r+4,0,!0),r+=8}),this._inst.exports.run(s,l),this.exited&&this._resolveExitPromise(),t.next=19,this._exitPromise;case 19:case"end":return t.stop()}},t,this)}));return function(e){return t.apply(this,arguments)}}()},{key:"_resume",value:function(){if(this.exited)throw new Error("Go program has already exited");this._inst.exports.resume(),this.exited&&this._resolveExitPromise()}},{key:"_makeFuncWrapper",value:function(e){var t=this;return function(){var n={id:e,this:this,args:arguments};return t._pendingEvent=n,t._resume(),n.result}}}]),t}(),o){t.argv.length<3&&(t.stderr.write("usage: go_js_wasm_exec [wasm binary] [arguments]\n"),t.exit(1));var p=new e.Go;p.argv=t.argv.slice(2),p.env=Object.assign({TMPDIR:n(602).tmpdir()},Object({NODE_ENV:"production",PUBLIC_URL:"/wdte"})),p.exit=t.exit,WebAssembly.instantiate(e.fs.readFileSync(t.argv[2]),p.importObject).then(function(e){return t.on("exit",function(e){0!==e||p.exited||(p._pendingEvent={id:0},p._resume())}),p.run(e.instance)}).catch(function(e){throw e})}}()}).call(this,n(24),n(25))},517:function(e,t){},519:function(e,t){},553:function(e,t){},554:function(e,t){},603:function(e,t,n){"use strict";n.r(t);var o={};n.r(o),n.d(o,"fib",function(){return R}),n.d(o,"stream",function(){return I}),n.d(o,"strings",function(){return N}),n.d(o,"lambdas",function(){return F}),n.d(o,"quine",function(){return A}),n.d(o,"hundredDoors",function(){return W});var i=n(0),r=n.n(i),a=n(38),s=n.n(a),c=(n(283),n(40)),u=n.n(c),l=n(66),d=n(68),m=n(150),f=n(29),p=n(625),h=n(623),g=n(20),v=n(243),w=n.n(v),b=n(244),x=n.n(b),y=n(105),E=n.n(y);n(374);E.a.define("ace/mode/wdte",function(e,t,n){var o=e("ace/lib/oop"),i=e("ace/mode/text").Mode,r=e("ace/mode/text_highlight_rules").TextHighlightRules,a=function(){this.$rules={start:[{token:"comment",regex:"^#.*$"},{token:"keyword",regex:"(\\b(memo|let|import)\\b)"},{token:"keyword.operator",regex:"\\.|\\{|\\}|\\[|\\]|\\(|\\)|=>|;|:|->|--|-\\||\\(@"},{token:"constant",regex:"\\b[0-9]+|[0-9]+\\.[0-9]+|\\.[0-9]+"},{token:"string",regex:"(([\"']).*\\2)"}]}};o.inherits(a,r);var s=function(){this.HighlightRules=a};o.inherits(s,i),t.Mode=s});var _=n(618),O=n(620),j=n(622),T=n(619),k=n(621),D=n(151),S=n.n(D),C="Introduction\n============\n\nWelcome to the WDTE playground, a browser based evaluation environment for WDTE. This playground's features includes the standard function set as well as a number of importable modules.\n\nIf you have never seen WDTE before and are completely confused at the moment, try reading the overview on the WDTE project's wiki: https://github.com/DeedleFake/wdte/wiki\n\nFun Fact\n--------\n\nThe WDTE interpreter has been compiled to WebAssembly for this playground, meaning that, by opening this page, you've downloaded the entire system. Congratulations.\n\nDocumentation\n-------------\n\nFor documentation on the standard function set, see https://godoc.org/github.com/DeedleFake/wdte/std\n\nThe standard library is available for importing, with the exception of the `io/file` module. The `io` module is pre-inserted into the initial scope as `io`. There is also a `playground` module which provides interaction with the playground. It is detailed below.\n\nPlayground Module\n-----------------\n\n#### wdteVersion\n    wdteVersion\nReturns the version of WDTE that the playground is using.\n\n#### goVersion\n    goVersion\nReturns the version of Go that the playground was built with.\n\nMacros\n------\n\n#### raw\nYields the raw text of its input.",U="io.stdout -> io.writeln 'Greetings, pocket universe.';",R={name:"Fibonacci",desc:"\nFibonacci\n=========\n\nThis example provides a memoized implementation of a recursive Fibonacci number generator. It also provides a recursive factorial function for the heck of it.\n",input:"let memo fib n => n {>= 2 => + (fib (- n 1)) (fib (- n 2))};\n\nlet ! n => n {\n\t<= 1 => 1;\n\ttrue => - n 1 -> ! -> * n;\n};\n\nfib 30\n-- io.writeln io.stdout\n-> / 5\n-- io.writeln io.stdout\n;"},I={name:"Stream",desc:"\nStream\n======\n\nThis example demonstrates the `stream` module. This module provides functional iterator operations, such as map, reduce, and filter.\n\nFor a full list of functions, see [the godocs][godoc].\n\n[godoc]: https://www.godoc.org/github.com/DeedleFake/wdte/std/stream\n",input:"let m => import 'math';\nlet s => import 'stream';\n\nio.writeln io.stdout 'Map and filter:';\ns.range 0 (* m.pi 2) (/ m.pi 2)\n-> s.map m.sin\n-> s.filter (>= 0)\n-> s.map (io.writeln io.stdout)\n-> s.drain\n;\n\nio.writeln io.stdout 'Reduce:';\ns.range 1 5\n-> s.reduce 1 *\n-- io.writeln io.stdout\n;"},N={name:"Strings",desc:"\nStrings\n=======\n\nThis example demonstrates the `strings` module. This module provides basic string operations, such as finding the index of a substring, as well as more complicated operations, such as formatting.\n\nFor a full list of functions, including an explanation of the formatting system, see [the godocs][godoc].\n\n[godoc]: https://www.godoc.org/github.com/DeedleFake/wdte/std/strings\n",input:"let a => import 'arrays';\nlet s => import 'stream';\nlet str => import 'strings';\n\na.stream ['abc'; 'bcd'; 'cde']\n-> s.map (str.index 'cd')\n-> s.collect\n-- io.writeln io.stdout\n;\n\n'This is the type of English up with which I will not put.'\n-> str.format '{q}'\n-- io.writeln io.stdout\n;"},F={name:"Lambdas",desc:"\nLambdas\n=======\n\nThis example demonstrates lambdas by implementing an iterative Fibonacci number calculator using streams.\n",input:"let s => import 'stream';\nlet a => import 'arrays';\n\nlet fib n => s.range 1 n\n\t-> s.reduce [0; 1] (@ self p n =>\n\t\tlet [a b] => p;\n\t\t[\n\t\t\tb;\n\t\t\t+ a b;\n\t\t];\n\t)\n\t-> at 1\n\t;\n\nfib 30\n-- io.writeln io.stdout\n;"},A={name:"Quine",desc:"\nQuine\n=====\n\nThis example is an implemenation of a quine. That's about it.\n",input:"let str => import 'strings';\nlet q => \"let str => import 'strings';\\nlet q => {q};\\nstr.format q q -- io.writeln io.stdout;\";\nstr.format q q -- io.writeln io.stdout;"},W={name:"100 Doors",desc:"\n100 Doors\n=========\n\nThe [100 doors problem](https://www.rosettacode.org/wiki/100_doors), as presented by Rosetta Code, is as follows:\n\nThere are 100 doors that are all closed. You walk past the doors in the same direction 100 times. On the first pass, you toggle the state of every door, opening closed doors and closing open doors. On the second pass, you toggle every second door. On the third you toggle every third door. Etc.\n\nThis example simulates this scenario, printing out the final state of the doors.\n",input:"let a => import 'arrays';\nlet s => import 'stream';\n\nlet toggle doors m =>\n\ta.stream doors\n\t-> s.enumerate\n\t-> s.map (@ s n => [+ (at n 0) 1; at n 1])\n\t-> s.map (@ s n => n {\n\t\t\t(@ s n => == (% (at n 0) m) 0) => ! (at n 1);\n\t\t\ttrue => at n 1;\n\t\t})\n\t-> s.collect\n\t;\n\ns.range 100\n-> s.map false\n-> s.collect : doors\n-> s.range 1 100\n-> s.reduce doors toggle\n-> a.stream\n-> s.map (@ s n => 0 {\n\t\tn => 'Open';\n\t\ttrue => 'Closed';\n\t} -- io.writeln io.stdout)\n-> s.drain\n;"},P=function(){var e=Object(l.a)(u.a.mark(function e(t){return u.a.wrap(function(e){for(;;)switch(e.prev=e.next){case 0:return e.next=2,new Promise(function(e,n){window.WDTE.run(t,function(t,o){if(null!=t)return n(t);e(o)})});case 2:return e.abrupt("return",e.sent);case 3:case"end":return e.stop()}},e)}));return function(t){return e.apply(this,arguments)}}(),M=function(e){var t=document.createElement("textarea");t.value=e;try{if(document.body.appendChild(t),t.focus(),t.select(),!document.execCommand("copy"))throw new Error("copy failed")}finally{document.body.removeChild(t)}},G=Object(_.a)(function(e){return{"@font-face":{fontFamily:"Go Mono",src:"url(assets/Go-Mono.ttf)"},main:{display:"flex",flexDirection:"row",backgroundColor:"#EEEEEE",boxSizing:"border-box",padding:8,width:"100%",height:"100%",position:"absolute"},column:{display:"flex",flexDirection:"column",flex:"1 0 300px",margin:8,overflowY:"auto"},message:{marginBottom:"8px !important"},inputToolbar:{flex:0},input:{flex:"0 1 50%",borderRadius:8},outputWrapper:{display:"flex",flexDirection:"column",flex:"0 1 50%",marginTop:12,boxShadow:"inset 4px 4px 4px #AAAAAA",borderRadius:8,backgroundColor:"#CCCCCC"},outputToolbar:{flex:0,alignSelf:"end",margin:"8px !important"},output:{flex:"1 0 0",fontFamily:"Go-Mono",fontSize:14,margin:"8px 8px 0px 8px",border:0,backgroundColor:"inherit",resize:"none",whiteSpace:"pre"},slide:{"&.enter":{transition:"all 300ms",overflow:"hidden",maxHeight:0,"&.active":{maxHeight:"500px"}},"&.exit":{transition:"all 300ms",overflow:"hidden",maxHeight:"500px","&.active":{maxHeight:0}}}}}),V=function(e){var t=G(),n=Object(i.useState)(C),a=Object(f.a)(n,2),s=a[0],c=a[1],v=Object(i.useState)(function(){try{return S.a.inflate(g.Buffer.from(window.location.hash.substr(1),"base64"),{to:"string"})}catch(e){return console.warn(e),U}}),b=Object(f.a)(v,2),y=b[0],E=b[1],_=Object(i.useState)(""),D=Object(f.a)(_,2),R=D[0],I=D[1],N=Object(i.useReducer)(function(e,t){switch(t.$){case"add":return null!=t.timeout&&setTimeout(function(){W({$:"remove",id:e.id})},t.timeout),Object(m.a)({},e,Object(d.a)({id:e.id+1},e.id,{type:t.type,msg:t.msg}));case"remove":return Object.entries(e).filter(function(e){var n=Object(f.a)(e,2),o=n[0];n[1];return o!==t.id.toString()}).reduce(function(e,t){var n=Object(f.a)(t,2),o=n[0],i=n[1];return Object(m.a)({},e,Object(d.a)({},o,i))},{});default:return e}},{id:0}),F=Object(f.a)(N,2),A=F[0],W=F[1],V=Object(i.useCallback)(function(e,t){var n=arguments.length>2&&void 0!==arguments[2]?arguments[2]:3e3;W({$:"add",type:e,msg:t,timeout:n})},[]),q=Object(i.useCallback)(Object(l.a)(u.a.mark(function e(){return u.a.wrap(function(e){for(;;)switch(e.prev=e.next){case 0:return e.prev=0,e.t0=I,e.next=4,P(y);case 4:e.t1=e.sent,(0,e.t0)(e.t1),e.next=11;break;case 8:e.prev=8,e.t2=e.catch(0),I(e.t2.toString());case 11:case"end":return e.stop()}},e,null,[[0,8]])})),[y]),L=Object(i.useCallback)(function(){try{var e=g.Buffer.from(S.a.deflate(y)).toString("base64");M("".concat(window.location.origin).concat(window.location.pathname,"#").concat(e)),window.location.href="#".concat(e),V("success","Link successfully copied to clipboard.")}catch(t){V("error","Failed to copy to clipboard: ".concat(t.toString()))}},[y,V]),H=Object(i.useCallback)(function(){try{M(R),V("success","Output successfully copied to clipboard.")}catch(e){V("error","Failed to copy to clipboard: ".concat(e.toString()))}},[R,V]);return r.a.createElement("div",{className:t.main},r.a.createElement(p.a,{component:"div",className:t.column},Object.entries(A).filter(function(e){var t=Object(f.a)(e,2),n=t[0];t[1];return!isNaN(n)}).map(function(e){var n=Object(f.a)(e,2),o=n[0],i=n[1];return r.a.createElement(h.a,{key:o,classNames:{enter:"enter",enterActive:"active",exit:"exit",exitActive:"active"},timeout:300},r.a.createElement(O.a,Object.assign({className:[t.message,t.slide].join(" ")},Object(d.a)({},i.type,!0)),r.a.createElement("p",null,i.msg)))}),r.a.createElement(w.a,{source:s})),r.a.createElement("div",{className:t.column},r.a.createElement(j.a,{className:t.inputToolbar,inverted:!0},r.a.createElement(j.a.Item,{onClick:q},"Run"),r.a.createElement(T.a,{item:!0,text:"Examples"},r.a.createElement(T.a.Menu,null,Object.entries(o).map(function(e){var t=Object(f.a)(e,2),n=t[0],i=t[1];return r.a.createElement(T.a.Item,{key:n,value:n,onClick:function(e,t){c(o[t.value].desc),E(o[t.value].input)}},i.name)}))),r.a.createElement(j.a.Item,{position:"right",onClick:L},"Share")),r.a.createElement(x.a,{className:t.input,style:{width:null,height:null},editorProps:{$blockScrolling:1/0},mode:"wdte",theme:"vibrant_ink",value:y,onChange:function(e){return E(e)}}),r.a.createElement("div",{className:t.outputWrapper},r.a.createElement("textarea",{className:t.output,readOnly:!0,value:R}),r.a.createElement(k.a.Group,{compact:!0,className:t.outputToolbar},r.a.createElement(k.a,{icon:"clipboard",onClick:H})))))},q=(n(514),new window.Go);WebAssembly.instantiateStreaming(fetch("./wdte.wasm"),q.importObject).then(function(e){return q.run(e.instance)}),s.a.render(r.a.createElement(V,null),document.getElementById("root"))}},[[278,1,2]]]);
//# sourceMappingURL=main.d63e3c59.chunk.js.map