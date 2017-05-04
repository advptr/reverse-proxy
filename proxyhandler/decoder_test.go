package proxyhandler

import (
	"strings"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewDecoder(t *testing.T) {
	assert := assert.New(t)

	schema, err := ParseXSDFile("../schema/MedicationSchema.xsd")
	assert.NoError(err)

	s := `<Envelope><Body><se:MedicationEHRExtracts xmlns:sch="http://www.ascc.net/xml/schematron"
xmlns:se="http://se.cambiosys.cosmicconnect.publish.xml.medication"
xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
xsi:schemaLocation="http://se.cambiosys.cosmicconnect.publish.xml.medication file:/Users/Peter/Documents/workspace/cambio/api/ccp/CCPKontrakt/ConnectPublishXML/XSD/MedicationSchema.xsd">
<XSDVersion>1.0</XSDVersion>
<MedicationEHRExtract>
<ehr_id/>
<ehr_system/>
<time_created>2017-03-07T12:00:00</time_created>
	<Medication>
	<rc_id/>
	<MedicationVersionDateTime>
	</MedicationVersionDateTime>
	<MedicationActive>false</MedicationActive>
	<MedicationStartDateTime>
	</MedicationStartDateTime>
	<MedicationChoice>
	</MedicationChoice>
	</Medication>
	<Patient>
	<PatientID>19121212-1212</PatientID>
	<PatientProtectedId>false</PatientProtectedId>
	<PatientSecrecyGroups>
	</PatientSecrecyGroups>
	<PatientName>
	<PatientFamilyName>PatientFamilyName0</PatientFamilyName>
	</PatientName>
	</Patient>
	<Units>
	<Unit>
	<UnitID/>
	<UnitName>UnitName0</UnitName>
	<UnitSecrecyGroups>
	<Group>A</Group>
	<Group>B</Group>
	</UnitSecrecyGroups>
	<AccessZones>
	</AccessZones>
	<PhoneNumbers>
	</PhoneNumbers>
	<EmailAddresses>
	</EmailAddresses>
	<UnitAddresses>
	</UnitAddresses>
	</Unit>
	</Units>
	<CreationUnit>
	<UnitID/>
	</CreationUnit>
	<ResponsibleUnit>
	<UnitID/>
	</ResponsibleUnit>
	<CareProviders>
	</CareProviders>
	</MedicationEHRExtract>
	<MedicationEHRExtract>
	<ehr_id/>
	<ehr_system/>
	<time_created>2017-04-09T14:00:00</time_created>
	<Medication>
	<rc_id/>
	<MedicationVersionDateTime>
	</MedicationVersionDateTime>
	<MedicationActive>false</MedicationActive>
	<MedicationStartDateTime>
	</MedicationStartDateTime>
	<MedicationChoice>
	</MedicationChoice>
	</Medication>
	<Patient>
	<PatientID/>
	<PatientProtectedId>false</PatientProtectedId>
	<PatientSecrecyGroups>
	</PatientSecrecyGroups>
	<PatientName>
	<PatientFamilyName>PatientFamilyName1</PatientFamilyName>
	</PatientName>
	</Patient>
	<Units>
	<Unit>
	<UnitID/>
	<UnitName>UnitName2</UnitName>
	<UnitSecrecyGroups>
	</UnitSecrecyGroups>
	<AccessZones>
	</AccessZones>
	<PhoneNumbers>
	</PhoneNumbers>
	<EmailAddresses>
	</EmailAddresses>
	<UnitAddresses>
	</UnitAddresses>
	</Unit>
	<Unit>
	<UnitID/>
	<UnitName>UnitName3</UnitName>
	<UnitSecrecyGroups>
	</UnitSecrecyGroups>
	<AccessZones>
	</AccessZones>
	<PhoneNumbers>
	</PhoneNumbers>
	<EmailAddresses>
	</EmailAddresses>
	<UnitAddresses>
	</UnitAddresses>
	</Unit>
	</Units>
	<CreationUnit>
	<UnitID/>
	</CreationUnit>
	<ResponsibleUnit>
	<UnitID/>
	</ResponsibleUnit>
	<CareProviders>
	</CareProviders>
	</MedicationEHRExtract>
	</se:MedicationEHRExtracts>
	</Body></Envelope>`

	
	schemaElements := make(map[string]xsdElement)
	for _, s := range schema {
		addElements(schemaElements,  s.Elements)
	}

	decoder := NewDecoder(strings.NewReader(s), schemaElements)

	root := &Node{}
	section, err := decoder.Decode(root, "MedicationEHRExtracts")
	assert.NoError(err)

	json, err := section.encode()
	assert.NoError(err)
	assert.NotNil(json)
}


func addElements(m map[string]xsdElement, elements []xsdElement) {
	for _, e := range elements {
		m[e.Name] = e
		if e.isComplex() {
			addElements(m, e.ComplexType.Sequence)
		}
	}
}