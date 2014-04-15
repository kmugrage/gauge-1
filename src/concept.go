package main

type conceptDictionary struct {
	conceptsMap map[string]*step
}

type conceptParser struct {

}

func (parser *conceptParser) parse(conceptString string) ([]*step,error) {
	return nil,nil
}

func (conceptDictionary * conceptDictionary) add(concepts []*step) {
	for _, step := range concepts {
		conceptDictionary.conceptsMap[step.value] = step
	}
}
