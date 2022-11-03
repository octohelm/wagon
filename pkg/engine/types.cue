pkg: {
	name:         string
	version:      string
	description?: string
}

// context: #Context & {}

dest: [Output=string]: {
	dest:   Output
	input:  #FS
	output: #FS
	...
}

#Step: {
	input:  #FS
	output: #FS
	...
}

//#Pipe: {
//	steps: [...#Step]
//
//	_dag: {
//		for idx, step in steps {
//			"\(idx)": step & {
//				if idx > 0 {
//					_desc:  _dag["\(idx-1)"].desc
//					source: _ | *_desc
//				}
//			}
//		}
//	}
//
//	if len(_dag) > 0 {
//		desc: _dag["\(len(_dag)-1)"].desc
//	}
//}
