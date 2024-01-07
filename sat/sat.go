package sat

import "fmt"

// Literal represents a variable with a boolean value
type Literal struct {
	Variable string
	Value    bool
}

// Clause is a disjunction of literals
type Clause []Literal

// Formula is a conjunction of clauses
type Formula []Clause

// Assignment maps variables to their boolean values
type Assignment map[string]bool

// SatisfyingASsignment finds a satisfying assignment for a given CNF formula
func SatisfyingAssignment(formula Formula, assignment Assignment) Assignment {
	if assignment == nil {
		assignment = make(Assignment)
	}
	// Unit propagation and simplification
	for {
		unitClauses := []Clause{}
		for _, clause := range formula {
			if len(clause) == 1 {
				unitClauses = append(unitClauses, clause)
			}
		}
		if len(unitClauses) == 0 {
			break
		}
		for _, clause := range unitClauses {
			variable, value := clause[0].Variable, clause[0].Value
			assignment[variable] = value
			var err error
			formula, err = ApplyAssignment(formula, variable, value)
			if err != nil {
				return nil
			}
		}
	}

	if len(formula) == 0 {
		return assignment
	}

	for _, clause := range formula {
		if len(clause) == 0 {
			return nil
		}
	}

	// Choose a variable to assign that is not already in the assignment
	var variable string
	for _, clause := range formula {
		for _, literal := range clause {
			if _, ok := assignment[literal.Variable]; !ok {
				variable = literal.Variable
				break
			}
		}
		if variable != "" {
			break
		}
	}

	for _, value := range []bool{true, false} {
		newAssignment := make(Assignment)
		for k, v := range assignment {
			newAssignment[k] = v
		}
		newAssignment[variable] = value
		newFormula, err := ApplyAssignment(formula, variable, value)
		if err != nil {
			continue
		}
		result := SatisfyingAssignment(newFormula, newAssignment)
		if result != nil {
			return result
		}
	}
	return nil
}

// ApplyAssignment applies an assignment to a formula
func ApplyAssignment(formula Formula, variable string, value bool) (Formula, error) {
	newFormula := Formula{}
	for _, clause := range formula {
		newClause := Clause{}
		for _, literal := range clause {
			if literal.Variable == variable {
				if literal.Value == value {
					// Clause is satisfied
					newClause = nil
					break
				}
				// Skip this literal as it makes the clause false
				continue
			}
			newClause = append(newClause, literal)
		}
		if newClause == nil {
			// Clause is satisfied, skip it
			continue
		}
		if len(newClause) == 0 {
			// Clause is unsatisfiable
			return nil, fmt.Errorf("unsatisfiable clause")
		}
		newFormula = append(newFormula, newClause)
	}

	return newFormula, nil
}
