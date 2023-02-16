package config

import (
	"fmt"
)

func (dredgeFile *DredgeFile) Validate() error {
	for _, r := range dredgeFile.Runtimes {
		if err := r.Validate(); err != nil {
			return err
		}
	}
	for _, w := range dredgeFile.Workflows {
		if err := w.Validate(); err != nil {
			return err
		}
	}
	for _, b := range dredgeFile.Buckets {
		if err := b.Validate(); err != nil {
			return err
		}
	}
	// TODO Validate resources here.
	return nil
}

func (r Runtime) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("name field is required for runtime")
	}
	if r.Type != RUNTIME_NATIVE && r.Type != RUNTIME_CONTAINER {
		return fmt.Errorf("unknown runtime type: %s (valid options are %s, %s)", r.Type, RUNTIME_NATIVE, RUNTIME_CONTAINER)
	}
	if r.Type == RUNTIME_NATIVE && (r.Image != "" ||
		r.Home != "" ||
		len(r.Cache) > 0 ||
		len(r.GlobalCache) > 0 ||
		len(r.Ports) > 0) {
		return fmt.Errorf("image, home, cache, global_cache and ports fields are only applicable to %s runtimes", RUNTIME_CONTAINER)
	}
	if r.Type == RUNTIME_CONTAINER && r.Image == "" {
		return fmt.Errorf("image field is required for %s runtimes", RUNTIME_CONTAINER)
	}
	return nil
}

func (b Bucket) Validate() error {
	if b.Name == "" {
		return fmt.Errorf("name field is required for bucket")
	}
	if b.Import != nil {
		if len(b.Workflows) > 0 {
			return fmt.Errorf("bucket %s: contains both workflows and an import", b.Name)
		}
		if err := b.Import.Validate(); err != nil {
			return fmt.Errorf("bucket %s: %v", b.Name, err)
		}
		return nil
	}
	for _, w := range b.Workflows {
		if err := w.Validate(); err != nil {
			return fmt.Errorf("bucket %s: %v", b.Name, err)
		}
	}
	return nil
}

func (i ImportBucket) Validate() error {
	if i.Bucket == "" {
		return fmt.Errorf("bucket field is required for import")
	}
	return nil
}

func (w Workflow) Validate() error {
	if w.Name == "" {
		return fmt.Errorf("name field is required for workflow")
	}
	if w.Import != nil {
		if len(w.Steps) > 0 {
			return fmt.Errorf("workflow %s: contains both steps and an import", w.Name)
		}
		if err := w.Import.Validate(); err != nil {
			return fmt.Errorf("workflow %s: %v", w.Name, err)
		}
		return nil
	}
	for _, i := range w.Inputs {
		if err := i.Validate(); err != nil {
			return fmt.Errorf("workflow %s: %v", w.Name, err)
		}
	}
	if len(w.Steps) == 0 {
		return fmt.Errorf("workflow %s: no steps or import defined", w.Name)
	}
	for _, s := range w.Steps {
		if err := s.Validate(); err != nil {
			return fmt.Errorf("workflow %s: %v", w.Name, err)
		}
	}
	return nil
}

func (i ImportWorkflow) Validate() error {
	if i.Workflow == "" {
		return fmt.Errorf("workflow field is required for import")
	}
	return nil
}

func (i Input) Validate() error {
	if i.Name == "" {
		return fmt.Errorf("name field is required on inputs")
	}
	if i.Type != "" && i.Type != INPUT_TEXT && i.Type != INPUT_SELECT {
		return fmt.Errorf("input %s: unknown input type: %s (valid options are: %s, %s)", i.Name, i.Type, INPUT_TEXT, INPUT_SELECT)
	}
	if len(i.Values) > 0 && (i.Type == "" || i.Type == INPUT_TEXT) {
		return fmt.Errorf("input %s: values for input can only be provided for the %s type", i.Name, INPUT_SELECT)
	}
	if len(i.Values) == 0 && i.Type == INPUT_SELECT {
		return fmt.Errorf("input %s: no values are provided, values are required for the %s type", i.Name, INPUT_SELECT)
	}
	if i.DefaultValue != "" && i.Type == INPUT_SELECT {
		return fmt.Errorf("input %s: default value can only be provided for the %s type", i.Name, INPUT_TEXT)
	}
	return nil
}

func (s Step) Validate() error {
	numFields := 0

	if s.Shell != nil {
		numFields += 1
		err := s.Shell.Validate()
		if err != nil {
			return err
		}
	}
	if s.Template != nil {
		numFields += 1
		err := s.Template.Validate()
		if err != nil {
			return err
		}
	}
	if s.Browser != nil {
		numFields += 1
		err := s.Browser.Validate()
		if err != nil {
			return err
		}
	}
	if s.EditDredgeFile != nil {
		numFields += 1
		err := s.EditDredgeFile.Validate()
		if err != nil {
			return err
		}
	}
	if s.If != nil {
		numFields += 1
		err := s.If.Validate()
		if err != nil {
			return err
		}
	}
	if s.Execute != nil {
		numFields += 1
		err := s.Execute.Validate()
		if err != nil {
			return err
		}
	}

	if numFields == 0 {
		return fmt.Errorf("step %s does not contain an action", s.Name)
	} else if numFields == 1 {
		return nil
	} else {
		return fmt.Errorf("step %s contains more than 1 action", s.Name)
	}
}

func (s ShellStep) Validate() error {
	if s.Cmd == "" {
		return fmt.Errorf("cmd field is required for shell")
	}
	return nil
}

func (t TemplateStep) Validate() error {
	if t.Input != "" && t.Source != "" {
		return fmt.Errorf("either input or source should be set for template")
	}
	if t.Dest == "" {
		return fmt.Errorf("dest field is required for template")
	}
	if t.Insert != nil {
		return t.Insert.Validate()
	}
	return nil
}

func (i Insert) Validate() error {
	if i.Placement != "" && i.Placement != INSERT_BEGIN && i.Placement != INSERT_END && i.Placement != INSERT_UNIQUE {
		return fmt.Errorf("unknown placement in insert: %s (valid options are: %s, %s, %s)", i.Placement, INSERT_BEGIN, INSERT_END, INSERT_UNIQUE)
	}
	return nil
}

func (b BrowserStep) Validate() error {
	if b.Url == "" {
		return fmt.Errorf("url field is required for browser")
	}
	return nil
}

func (e EditDredgeFileStep) Validate() error {
	return nil
}

func (i IfStep) Validate() error {
	if i.Cond == "" {
		return fmt.Errorf("cond field is required for if")
	}
	if len(i.Steps) == 0 {
		return fmt.Errorf("1 or more steps are required for if")
	}
	return nil
}

func (e ExecuteStep) Validate() error {
	if e.Resource == "" {
		return fmt.Errorf("resource field is required for execute")
	}
	if e.Command == "" {
		return fmt.Errorf("command field is required for execute")
	}
	return nil
}
