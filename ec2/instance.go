package ec2


import (
	"xml"
)
/*
  The XML package doesn't do deeply-nested types with reflection well,
 so much of the functionality provided is simply hacks
 around the very-deep responses given by EC2.
*/

type describeInstancesResponse struct {
	RequestId    string           "requestId"
	Reservations []ReservationSet "reservationSet>item"
}

type GroupSet struct {
	GroupId   string "groupId"
	GroupName string "groupName"
}

type ReservationSet struct {
	ReservationId string
	OwnerId       string
	Groups        []GroupSet "groupSet>item"
	Instances     []Instance "instancesSet>item"
}


// At the level of depth we're at, XML's 
//  not happy with us
type Instance struct {
	XMLName          xml.Name
	InstanceId       string
	ImageId          string
	PrivateDNSName   string
	DNSName          string
	PrivateIPAddress string
	IPAddress        string
	AvailabilityZone string "placement>availabilityZone"
	MonitoringState  string "monitoring>state"
	InstanceType     string
	RootDeviceName   string
	RootDeviceType   string
	KernelId         string
}

/* Example EC2 Output (2011-06-21)

<DescribeInstancesResponse xmlns="http://ec2.amazonaws.com/doc/2011-05-15/">
    <requestId>9e7685a5-dc96-4ee8-b587-0d4c3b42aed6</requestId>
    <reservationSet>
        <item>
            <reservationId>r-8982afe5</reservationId>
            <ownerId>930374178234</ownerId>
            <groupSet>
                <item>
                    <groupId>sg-b0dc65d9</groupId>
                    <groupName>default</groupName>
                </item>
            </groupSet>
            <instancesSet>
                <item>
                    <instanceId>i-21672e4f</instanceId>
                    <imageId>ami-e00df089</imageId>
                    <instanceState>
                        <code>16</code>
                        <name>running</name>
                    </instanceState>
                    <privateDnsName>domU-12-31-39-0B-0E-8E.compute-1.internal</privateDnsName>
                    <dnsName>ec2-184-72-81-29.compute-1.amazonaws.com</dnsName>
                    <reason/>
                    <keyName>oddgeneration</keyName>
                    <amiLaunchIndex>0</amiLaunchIndex>
                    <productCodes/>
                    <instanceType>t1.micro</instanceType>
                    <launchTime>2011-06-18T17:19:07.000Z</launchTime>
                    <placement>
                        <availabilityZone>us-east-1a</availabilityZone>
                        <groupName/>
                        <tenancy>default</tenancy>
                    </placement>
                    <kernelId>aki-4e7d9527</kernelId>
                    <monitoring>
                        <state>disabled</state>
                    </monitoring>
                    <privateIpAddress>10.214.17.120</privateIpAddress>
                    <ipAddress>184.72.81.29</ipAddress>
                    <groupSet>
                        <item>
                            <groupId>sg-b0dc65d9</groupId>
                            <groupName>default</groupName>
                        </item>
                    </groupSet>
                    <architecture>x86_64</architecture>
                    <rootDeviceType>ebs</rootDeviceType>
                    <rootDeviceName>/dev/sda1</rootDeviceName>
                    <blockDeviceMapping>
                        <item>
                            <deviceName>/dev/sda</deviceName>
                            <ebs>
                                <volumeId>vol-fbf98690</volumeId>
                                <status>attached</status>
                                <attachTime>2011-06-18T17:19:30.000Z</attachTime>
                                <deleteOnTermination>true</deleteOnTermination>
                            </ebs>
                        </item>
                    </blockDeviceMapping>
                    <virtualizationType>paravirtual</virtualizationType>
                    <clientToken/>
                    <tagSet>
                        <item>
                            <key>Name</key>
                            <value>build0</value>
                        </item>
                    </tagSet>
                    <hypervisor>xen</hypervisor>
                </item>
            </instancesSet>
            <requesterId>058890971305</requesterId>
        </item>
        <item>
            <reservationId>r-5fc2eb33</reservationId>
            <ownerId>930374178234</ownerId>
            <groupSet>
                <item>
                    <groupId>sg-b0dc65d9</groupId>
                    <groupName>default</groupName>
                </item>
            </groupSet>
            <instancesSet>
                <item>
                    <instanceId>i-8782d6e9</instanceId>
                    <imageId>ami-e00df089</imageId>
                    <instanceState>
                        <code>16</code>
                        <name>running</name>
                    </instanceState>
                    <privateDnsName>ip-10-122-46-4.ec2.internal</privateDnsName>
                    <dnsName>ec2-50-17-140-146.compute-1.amazonaws.com</dnsName>
                    <reason/>
                    <keyName>oddgeneration</keyName>
                    <amiLaunchIndex>0</amiLaunchIndex>
                    <productCodes/>
                    <instanceType>t1.micro</instanceType>
                    <launchTime>2011-06-20T05:53:35.000Z</launchTime>
                    <placement>
                        <availabilityZone>us-east-1c</availabilityZone>
                        <groupName/>
                        <tenancy>default</tenancy>
                    </placement>
                    <kernelId>aki-4e7d9527</kernelId>
                    <monitoring>
                        <state>disabled</state>
                    </monitoring>
                    <privateIpAddress>10.122.46.4</privateIpAddress>
                    <ipAddress>50.17.140.146</ipAddress>
                    <groupSet>
                        <item>
                            <groupId>sg-b0dc65d9</groupId>
                            <groupName>default</groupName>
                        </item>
                    </groupSet>
                    <architecture>x86_64</architecture>
                    <rootDeviceType>ebs</rootDeviceType>
                    <rootDeviceName>/dev/sda1</rootDeviceName>
                    <blockDeviceMapping>
                        <item>
                            <deviceName>/dev/sda</deviceName>
                            <ebs>
                                <volumeId>vol-dbbfc7b0</volumeId>
                                <status>attached</status>
                                <attachTime>2011-06-20T05:53:54.000Z</attachTime>
                                <deleteOnTermination>true</deleteOnTermination>
                            </ebs>
                        </item>
                    </blockDeviceMapping>
                    <instanceLifecycle>spot</instanceLifecycle>
                    <spotInstanceRequestId>sir-fa163214</spotInstanceRequestId>
                    <virtualizationType>paravirtual</virtualizationType>
                    <clientToken/>
                    <hypervisor>xen</hypervisor>
                </item>
            </instancesSet>
            <requesterId>854251627541</requesterId>
        </item>
    </reservationSet>
</DescribeInstancesResponse>
*/
