package eks

const (
	serviceRoleTemplate = `
---
AWSTemplateFormatVersion: '2010-09-09'
Description: 'Amazon EKS Service Role'


Resources:

  AWSServiceRoleForAmazonEKS:
    Type: AWS::IAM::Role
    Properties:
      AssumeRolePolicyDocument:
        Version: '2012-10-17'
        Statement:
        - Effect: Allow
          Principal:
            Service:
            - eks.amazonaws.com
          Action:
          - sts:AssumeRole
      ManagedPolicyArns:
        - arn:aws:iam::aws:policy/AmazonEKSServicePolicy
        - arn:aws:iam::aws:policy/AmazonEKSClusterPolicy

Outputs:

  RoleArn:
    Description: The role that EKS will use to create AWS resources for Kubernetes clusters
    Value: !GetAtt AWSServiceRoleForAmazonEKS.Arn
    Export:
      Name: !Sub "${AWS::StackName}-RoleArn"
`
	vpcTemplate = `
---
AWSTemplateFormatVersion: '2010-09-09'
Description: 'Amazon EKS Sample VPC'

Parameters:

  VpcBlock:
    Type: String
    Default: 192.168.0.0/16
    Description: The CIDR range for the VPC. This should be a valid private (RFC 1918) CIDR range.

  Subnet01Block:
    Type: String
    Default: 192.168.64.0/18
    Description: CidrBlock for subnet 01 within the VPC

  Subnet02Block:
    Type: String
    Default: 192.168.128.0/18
    Description: CidrBlock for subnet 02 within the VPC

  Subnet03Block:
    Type: String
    Default: 192.168.192.0/18
    Description: CidrBlock for subnet 03 within the VPC

Metadata:
  AWS::CloudFormation::Interface:
    ParameterGroups:
      -
        Label:
          default: "Worker Network Configuration"
        Parameters:
          - VpcBlock
          - Subnet01Block
          - Subnet02Block
          - Subnet03Block

Resources:
  VPC:
    Type: AWS::EC2::VPC
    Properties:
      CidrBlock:  !Ref VpcBlock
      EnableDnsSupport: true
      EnableDnsHostnames: true
      Tags:
      - Key: Name
        Value: !Sub '${AWS::StackName}-VPC'

  InternetGateway:
    Type: "AWS::EC2::InternetGateway"

  VPCGatewayAttachment:
    Type: "AWS::EC2::VPCGatewayAttachment"
    Properties:
      InternetGatewayId: !Ref InternetGateway
      VpcId: !Ref VPC

  RouteTable:
    Type: AWS::EC2::RouteTable
    Properties:
      VpcId: !Ref VPC
      Tags:
      - Key: Name
        Value: Public Subnets
      - Key: Network
        Value: Public

  Route:
    DependsOn: VPCGatewayAttachment
    Type: AWS::EC2::Route
    Properties:
      RouteTableId: !Ref RouteTable
      DestinationCidrBlock: 0.0.0.0/0
      GatewayId: !Ref InternetGateway

  Subnet01:
    Type: AWS::EC2::Subnet
    Metadata:
      Comment: Subnet 01
    Properties:
      AvailabilityZone:
        Fn::Select:
        - '0'
        - Fn::GetAZs:
            Ref: AWS::Region
      CidrBlock:
        Ref: Subnet01Block
      VpcId:
        Ref: VPC
      Tags:
      - Key: Name
        Value: !Sub "${AWS::StackName}-Subnet01"

  Subnet02:
    Type: AWS::EC2::Subnet
    Metadata:
      Comment: Subnet 02
    Properties:
      AvailabilityZone:
        Fn::Select:
        - '1'
        - Fn::GetAZs:
            Ref: AWS::Region
      CidrBlock:
        Ref: Subnet02Block
      VpcId:
        Ref: VPC
      Tags:
      - Key: Name
        Value: !Sub "${AWS::StackName}-Subnet02"

  Subnet03:
    Type: AWS::EC2::Subnet
    Metadata:
      Comment: Subnet 03
    Properties:
      AvailabilityZone:
        Fn::Select:
        - '2'
        - Fn::GetAZs:
            Ref: AWS::Region
      CidrBlock:
        Ref: Subnet03Block
      VpcId:
        Ref: VPC
      Tags:
      - Key: Name
        Value: !Sub "${AWS::StackName}-Subnet03"

  Subnet01RouteTableAssociation:
    Type: AWS::EC2::SubnetRouteTableAssociation
    Properties:
      SubnetId: !Ref Subnet01
      RouteTableId: !Ref RouteTable

  Subnet02RouteTableAssociation:
    Type: AWS::EC2::SubnetRouteTableAssociation
    Properties:
      SubnetId: !Ref Subnet02
      RouteTableId: !Ref RouteTable

  Subnet03RouteTableAssociation:
    Type: AWS::EC2::SubnetRouteTableAssociation
    Properties:
      SubnetId: !Ref Subnet03
      RouteTableId: !Ref RouteTable

  ControlPlaneSecurityGroup:
    Type: AWS::EC2::SecurityGroup
    Properties:
      GroupDescription: Cluster communication with worker nodes
      VpcId: !Ref VPC

Outputs:

  SubnetIds:
    Description: All subnets in the VPC
    Value: !Join [ ",", [ !Ref Subnet01, !Ref Subnet02, !Ref Subnet03 ] ]

  SecurityGroups:
    Description: Security group for the cluster control plane communication with worker nodes
    Value: !Join [ ",", [ !Ref ControlPlaneSecurityGroup ] ]

  VpcId:
    Description: The VPC Id
    Value: !Ref VPC
`
	workerNodesTemplate = `
---
AWSTemplateFormatVersion: '2010-09-09'
Description: 'Amazon EKS - Node Group'

Parameters:

  KeyName:
    Description: The EC2 Key Pair to allow SSH access to the instances
    Type: AWS::EC2::KeyPair::KeyName

  NodeImageId:
    Type: AWS::EC2::Image::Id
    Description: AMI id for the node instances.

  NodeInstanceType:
    Description: EC2 instance type for the node instances
    Type: String
    Default: t2.medium
    AllowedValues:
    - t2.small
    - t2.medium
    - t2.large
    - t2.xlarge
    - t2.2xlarge
    - m3.medium
    - m3.large
    - m3.xlarge
    - m3.2xlarge
    - m4.large
    - m4.xlarge
    - m4.2xlarge
    - m4.4xlarge
    - m4.10xlarge
    - m5.large
    - m5.xlarge
    - m5.2xlarge
    - m5.4xlarge
    - m5.12xlarge
    - m5.24xlarge
    - c4.large
    - c4.xlarge
    - c4.2xlarge
    - c4.4xlarge
    - c4.8xlarge
    - c5.large
    - c5.xlarge
    - c5.2xlarge
    - c5.4xlarge
    - c5.9xlarge
    - c5.18xlarge
    - i3.large
    - i3.xlarge
    - i3.2xlarge
    - i3.4xlarge
    - i3.8xlarge
    - i3.16xlarge
    - r3.xlarge
    - r3.2xlarge
    - r3.4xlarge
    - r3.8xlarge
    - r4.large
    - r4.xlarge
    - r4.2xlarge
    - r4.4xlarge
    - r4.8xlarge
    - r4.16xlarge
    - x1.16xlarge
    - x1.32xlarge
    - p2.xlarge
    - p2.8xlarge
    - p2.16xlarge
    - p3.2xlarge
    - p3.8xlarge
    - p3.16xlarge
    ConstraintDescription: must be a valid EC2 instance type

  NodeAutoScalingGroupMinSize:
    Type: Number
    Description: Minimum size of Node Group ASG.
    Default: 1

  NodeAutoScalingGroupMaxSize:
    Type: Number
    Description: Maximum size of Node Group ASG.
    Default: 3

  ClusterName:
    Description: The cluster name provided when the cluster was created.  If it is incorrect, nodes will not be able to join the cluster.
    Type: String

  NodeGroupName:
    Description: Unique identifier for the Node Group.
    Type: String

  ClusterControlPlaneSecurityGroup:
    Description: The security group of the cluster control plane.
    Type: AWS::EC2::SecurityGroup::Id

  VpcId:
    Description: The VPC of the worker instances
    Type: AWS::EC2::VPC::Id

  Subnets:
    Description: The subnets where workers can be created.
    Type: List<AWS::EC2::Subnet::Id>

Mappings:
  MaxPodsPerNode:
    c4.large:
      MaxPods: 29
    c4.xlarge:
      MaxPods: 58
    c4.2xlarge:
      MaxPods: 58
    c4.4xlarge:
      MaxPods: 234
    c4.8xlarge:
      MaxPods: 234
    c5.large:
      MaxPods: 29
    c5.xlarge:
      MaxPods: 58
    c5.2xlarge:
      MaxPods: 58
    c5.4xlarge:
      MaxPods: 234
    c5.9xlarge:
      MaxPods: 234
    c5.18xlarge:
      MaxPods: 737
    i3.large:
      MaxPods: 29
    i3.xlarge:
      MaxPods: 58
    i3.2xlarge:
      MaxPods: 58
    i3.4xlarge:
      MaxPods: 234
    i3.8xlarge:
      MaxPods: 234
    i3.16xlarge:
      MaxPods: 737
    m3.medium:
      MaxPods: 12
    m3.large:
      MaxPods: 29
    m3.xlarge:
      MaxPods: 58
    m3.2xlarge:
      MaxPods: 118
    m4.large:
      MaxPods: 20
    m4.xlarge:
      MaxPods: 58
    m4.2xlarge:
      MaxPods: 58
    m4.4xlarge:
      MaxPods: 234
    m4.10xlarge:
      MaxPods: 234
    m5.large:
      MaxPods: 29
    m5.xlarge:
      MaxPods: 58
    m5.2xlarge:
      MaxPods: 58
    m5.4xlarge:
      MaxPods: 234
    m5.12xlarge:
      MaxPods: 234
    m5.24xlarge:
      MaxPods: 737
    p2.xlarge:
      MaxPods: 58
    p2.8xlarge:
      MaxPods: 234
    p2.16xlarge:
      MaxPods: 234
    p3.2xlarge:
      MaxPods: 58
    p3.8xlarge:
      MaxPods: 234
    p3.16xlarge:
      MaxPods: 234
    r3.xlarge:
      MaxPods: 58
    r3.2xlarge:
      MaxPods: 58
    r3.4xlarge:
      MaxPods: 234
    r3.8xlarge:
      MaxPods: 234
    r4.large:
      MaxPods: 29
    r4.xlarge:
      MaxPods: 58
    r4.2xlarge:
      MaxPods: 58
    r4.4xlarge:
      MaxPods: 234
    r4.8xlarge:
      MaxPods: 234
    r4.16xlarge:
      MaxPods: 737
    t2.small:
      MaxPods: 8
    t2.medium:
      MaxPods: 17
    t2.large:
      MaxPods: 35
    t2.xlarge:
      MaxPods: 44
    t2.2xlarge:
      MaxPods: 44
    x1.16xlarge:
      MaxPods: 234
    x1.32xlarge:
      MaxPods: 234

Metadata:
  AWS::CloudFormation::Interface:
    ParameterGroups:
      -
        Label:
          default: "EKS Cluster"
        Parameters:
          - ClusterName
          - ClusterControlPlaneSecurityGroup
      -
        Label:
          default: "Worker Node Configuration"
        Parameters:
          - NodeGroupName
          - NodeAutoScalingGroupMinSize
          - NodeAutoScalingGroupMaxSize
          - NodeInstanceType
          - NodeImageId
          - KeyName
      -
        Label:
          default: "Worker Network Configuration"
        Parameters:
          - VpcId
          - Subnets

Resources:

  NodeInstanceProfile:
    Type: AWS::IAM::InstanceProfile
    Properties:
      Path: "/"
      Roles:
      - !Ref NodeInstanceRole

  NodeInstanceRole:
    Type: AWS::IAM::Role
    Properties:
      AssumeRolePolicyDocument:
        Version: '2012-10-17'
        Statement:
        - Effect: Allow
          Principal:
            Service:
            - ec2.amazonaws.com
          Action:
          - sts:AssumeRole
      Path: "/"
      ManagedPolicyArns:
        - arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy
        - arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy
        - arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly

  NodeSecurityGroup:
    Type: AWS::EC2::SecurityGroup
    Properties:
      GroupDescription: Security group for all nodes in the cluster
      VpcId:
        !Ref VpcId
      Tags:
      - Key: !Sub "kubernetes.io/cluster/${ClusterName}"
        Value: 'owned'

  NodeSecurityGroupIngress:
    Type: AWS::EC2::SecurityGroupIngress
    DependsOn: NodeSecurityGroup
    Properties:
      Description: Allow node to communicate with each other
      GroupId: !Ref NodeSecurityGroup
      SourceSecurityGroupId: !Ref NodeSecurityGroup
      IpProtocol: '-1'
      FromPort: 0
      ToPort: 65535

  NodeSecurityGroupFromControlPlaneIngress:
    Type: AWS::EC2::SecurityGroupIngress
    DependsOn: NodeSecurityGroup
    Properties:
      Description: Allow worker Kubelets and pods to receive communication from the cluster control plane
      GroupId: !Ref NodeSecurityGroup
      SourceSecurityGroupId: !Ref ClusterControlPlaneSecurityGroup
      IpProtocol: tcp
      FromPort: 1025
      ToPort: 65535

  ControlPlaneEgressToNodeSecurityGroup:
    Type: AWS::EC2::SecurityGroupEgress
    DependsOn: NodeSecurityGroup
    Properties:
      Description: Allow the cluster control plane to communicate with worker Kubelet and pods
      GroupId: !Ref ClusterControlPlaneSecurityGroup
      DestinationSecurityGroupId: !Ref NodeSecurityGroup
      IpProtocol: tcp
      FromPort: 1025
      ToPort: 65535

  ClusterControlPlaneSecurityGroupIngress:
    Type: AWS::EC2::SecurityGroupIngress
    DependsOn: NodeSecurityGroup
    Properties:
      Description: Allow pods to communicate with the cluster API Server
      GroupId: !Ref ClusterControlPlaneSecurityGroup
      SourceSecurityGroupId: !Ref NodeSecurityGroup
      IpProtocol: tcp
      ToPort: 443
      FromPort: 443

  NodeGroup:
    Type: AWS::AutoScaling::AutoScalingGroup
    Properties:
      DesiredCapacity: !Ref NodeAutoScalingGroupMaxSize
      LaunchConfigurationName: !Ref NodeLaunchConfig
      MinSize: !Ref NodeAutoScalingGroupMinSize
      MaxSize: !Ref NodeAutoScalingGroupMaxSize
      VPCZoneIdentifier:
        !Ref Subnets
      Tags:
      - Key: Name
        Value: !Sub "${ClusterName}-${NodeGroupName}-Node"
        PropagateAtLaunch: 'true'
      - Key: !Sub 'kubernetes.io/cluster/${ClusterName}'
        Value: 'owned'
        PropagateAtLaunch: 'true'
    UpdatePolicy:
      AutoScalingRollingUpdate:
        MinInstancesInService: '1'
        MaxBatchSize: '1'

  NodeLaunchConfig:
    Type: AWS::AutoScaling::LaunchConfiguration
    Properties:
      AssociatePublicIpAddress: 'true'
      IamInstanceProfile: !Ref NodeInstanceProfile
      ImageId: !Ref NodeImageId
      InstanceType: !Ref NodeInstanceType
      KeyName: !Ref KeyName
      SecurityGroups:
      - !Ref NodeSecurityGroup
      UserData:
        Fn::Base64:
          Fn::Join: [
            "",
            [
              "#!/bin/bash -xe\n",
              "CA_CERTIFICATE_DIRECTORY=/etc/kubernetes/pki", "\n",
              "CA_CERTIFICATE_FILE_PATH=$CA_CERTIFICATE_DIRECTORY/ca.crt", "\n",
              "MODEL_DIRECTORY_PATH=~/.aws/eks", "\n",
              "MODEL_FILE_PATH=$MODEL_DIRECTORY_PATH/eks-2017-11-01.normal.json", "\n",
              "mkdir -p $CA_CERTIFICATE_DIRECTORY", "\n",
              "mkdir -p $MODEL_DIRECTORY_PATH", "\n",
              "curl -o $MODEL_FILE_PATH https://s3-us-west-2.amazonaws.com/amazon-eks/1.10.3/2018-06-05/eks-2017-11-01.normal.json", "\n",
              "aws configure add-model --service-model file://$MODEL_FILE_PATH --service-name eks", "\n",
              "aws eks describe-cluster --region=", { Ref: "AWS::Region" }," --name=", { Ref: ClusterName }," --query 'cluster.{certificateAuthorityData: certificateAuthority.data, endpoint: endpoint}' > /tmp/describe_cluster_result.json", "\n",
              "cat /tmp/describe_cluster_result.json | grep certificateAuthorityData | awk '{print $2}' | sed 's/[,\"]//g' | base64 -d >  $CA_CERTIFICATE_FILE_PATH", "\n",
              "MASTER_ENDPOINT=$(cat /tmp/describe_cluster_result.json | grep endpoint | awk '{print $2}' | sed 's/[,\"]//g')", "\n",
              "INTERNAL_IP=$(curl -s http://169.254.169.254/latest/meta-data/local-ipv4)", "\n",
              "sed -i s,MASTER_ENDPOINT,$MASTER_ENDPOINT,g /var/lib/kubelet/kubeconfig", "\n",
              "sed -i s,CLUSTER_NAME,", { Ref: ClusterName }, ",g /var/lib/kubelet/kubeconfig", "\n",
              "sed -i s,REGION,", { Ref: "AWS::Region" }, ",g /etc/systemd/system/kubelet.service", "\n",
              "sed -i s,MAX_PODS,", { "Fn::FindInMap": [ MaxPodsPerNode, { Ref: NodeInstanceType }, MaxPods ] }, ",g /etc/systemd/system/kubelet.service", "\n",
              "sed -i s,MASTER_ENDPOINT,$MASTER_ENDPOINT,g /etc/systemd/system/kubelet.service", "\n",
              "sed -i s,INTERNAL_IP,$INTERNAL_IP,g /etc/systemd/system/kubelet.service", "\n",
              "DNS_CLUSTER_IP=10.100.0.10", "\n",
              "if [[ $INTERNAL_IP == 10.* ]] ; then DNS_CLUSTER_IP=172.20.0.10; fi", "\n",
              "sed -i s,DNS_CLUSTER_IP,$DNS_CLUSTER_IP,g  /etc/systemd/system/kubelet.service", "\n",
              "sed -i s,CERTIFICATE_AUTHORITY_FILE,$CA_CERTIFICATE_FILE_PATH,g /var/lib/kubelet/kubeconfig" , "\n",
              "sed -i s,CLIENT_CA_FILE,$CA_CERTIFICATE_FILE_PATH,g  /etc/systemd/system/kubelet.service" , "\n",
              "systemctl daemon-reload", "\n",
              "systemctl restart kubelet", "\n",
              "/opt/aws/bin/cfn-signal -e $? ",
              "         --stack ", { Ref: "AWS::StackName" },
              "         --resource NodeGroup ",
              "         --region ", { Ref: "AWS::Region" }, "\n"
            ]
          ]

Outputs:
  NodeInstanceRole:
    Description: The node instance role
    Value: !GetAtt NodeInstanceRole.Arn
`
)
