package app

func (w *fileWork) xRcvWithSelect() {
	name_x := w.entity.Name()
	x_get := "*" + name_x + "GetRequest"
	w.P("func (x ", x_get, ") WithSelect(f func(s *", name_x+"Select", ")) ", x_get, " {")
	w.P("	if !x.HasSelect() {")
	w.P("		x.SetSelect(&", name_x+"Select", "{})")
	w.P("	}")
	w.P("	f(x.GetSelect())")
	w.P("	return x")
	w.P("}")
	w.P("")
}
