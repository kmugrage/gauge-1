package main

type conceptDictionary struct {
	conceptsMap map[string]*step
}

type conceptParser struct {

}

func (parser *conceptParser) parse(text string) ([]*step, error) {
//	specParser := new(specParser)
//	tokens, err := specParser.generateTokens(text)
//	if (err != nil) {
//		return nil, err
//	}
	return nil, nil
}

func (conceptDictionary * conceptDictionary) add(concepts []*step) {
	for _, step := range concepts {
		conceptDictionary.conceptsMap[step.value] = step
	}
}
