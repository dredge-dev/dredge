package config

import (
	"fmt"
)

func (dredgeFile *DredgeFile) Validate() error {
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
