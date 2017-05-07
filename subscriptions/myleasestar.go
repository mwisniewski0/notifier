package subscriptions

import (
  "net/http"
  "io/ioutil"
  "fmt"
  "encoding/json"
)

type MyLeaseStarVacancy struct {
  availableUnitsList []int
  ComplexName string
  PropertyId string
}

type UnitInformation struct {
  // Using floats, because the IP returns them sometimes
  Rent float64
  NumberOfBaths float64
  NumberOfBeds float64
  InternalAvailableDate string
  VacantDate string
  Id float64
  SquareFeet float64
}

type myLeaseStarVacancyResponse struct {
  Units []UnitInformation
}

func (checker *MyLeaseStarVacancy) NextCheck() (*string, error) {
  response, err := http.Get("http://api.myleasestar.com/v2/property/" +
    checker.PropertyId + "/units?available=true")
  if err != nil {
	  fmt.Println("A request to MyLeaseStar failed")
    return nil, err
  }

  defer response.Body.Close()

  // parse the Json response
  body, err := ioutil.ReadAll(response.Body)
  if err != nil {
    fmt.Println("Could not read response body")
    return nil, err
  }
  parsedResponse := new(myLeaseStarVacancyResponse)
  err = json.Unmarshal(body, parsedResponse)
  if err != nil {
    fmt.Println("Response from MyLeaseStart was invalid")
    return nil, err
  }

  newVacantUnits := make([]UnitInformation, 0) // Wil contain all newly vacant units
  updatedAvailableUnitsList := make([]int, 0) // Will replace the old list of available units
  for _, unit := range(parsedResponse.Units) {
    updatedAvailableUnitsList = append(updatedAvailableUnitsList, int(unit.Id))

    // check if the unit was available previously
    unitAvailablePreviously := false
    for _, previousId := range(checker.availableUnitsList) {
      if int(unit.Id) == previousId {
        unitAvailablePreviously = true
        break
      }
    }

    if !unitAvailablePreviously {
      // Hooray! A new unit has been found!
      newVacantUnits = append(newVacantUnits, unit)
    }
  }

  checker.availableUnitsList = updatedAvailableUnitsList
  if len(newVacantUnits) == 0{
    // nothing to report
    return nil, nil
  } else {
    message := "Subject: New vacancies at " + checker.ComplexName + "!\n\n"
    message += "New vacancies:\n\n"
    jsonEncoded, err := json.MarshalIndent(newVacantUnits, "", "  ")
    if err != nil {
      fmt.Println("Could not covert the list of vacant units to JSON")
      return nil, err
    }
    message += string(jsonEncoded)
    return &message, nil
  }
}

func (checker *MyLeaseStarVacancy) IsConcurrent() bool {
  return false
}
