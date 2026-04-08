import type { Rocket, Launch, News, LaunchBase, Company } from '../types';

export const rockets: Rocket[] = [
  {
    id: 1,
    name: 'Falcon 9',
    description: 'A reusable, two-stage rocket designed and manufactured by SpaceX for the reliable and safe transport of people and payloads into Earth orbit and beyond.',
    height: 70,
    diameter: 3.7,
    mass: 549054,
    company: 'SpaceX',
    imageUrl: 'https://via.placeholder.com/400x600/1e3a8a/ffffff?text=Falcon+9',
    active: true,
  },
  {
    id: 2,
    name: 'Falcon Heavy',
    description: 'The most powerful operational rocket in the world by a factor of two, capable of carrying large payloads to orbit and beyond.',
    height: 70,
    diameter: 12.2,
    mass: 1420788,
    company: 'SpaceX',
    imageUrl: 'https://via.placeholder.com/400x600/7c3aed/ffffff?text=Falcon+Heavy',
    active: true,
  },
  {
    id: 3,
    name: 'Starship',
    description: 'A fully reusable transportation system designed to carry both crew and cargo to Earth orbit, the Moon, Mars, and beyond.',
    height: 120,
    diameter: 9,
    mass: 5000000,
    company: 'SpaceX',
    imageUrl: 'https://via.placeholder.com/400x600/059669/ffffff?text=Starship',
    active: true,
  },
  {
    id: 4,
    name: 'New Shepard',
    description: 'A fully reusable sub-orbital launch vehicle developed by Blue Origin for space tourism.',
    height: 18,
    diameter: 3.7,
    mass: 75000,
    company: 'Blue Origin',
    imageUrl: 'https://via.placeholder.com/400x600/dc2626/ffffff?text=New+Shepard',
    active: true,
  },
];

export const launches: Launch[] = [
  {
    id: 1,
    name: 'Starlink Mission',
    date: '2025-10-23T10:30:00Z',
    rocket: 'Falcon 9',
    launchBase: 'Kennedy Space Center',
    status: 'scheduled',
    description: 'Deployment of 60 Starlink satellites to low Earth orbit.',
  },
  {
    id: 2,
    name: 'GPS III SV06',
    date: '2025-10-25T09:15:00Z',
    rocket: 'Falcon 9',
    launchBase: 'Cape Canaveral',
    status: 'scheduled',
    description: 'Launch of GPS III satellite for the US Space Force.',
  },
  {
    id: 3,
    name: 'OneWeb Mission',
    date: '2025-10-27T14:20:00Z',
    rocket: 'Falcon 9',
    launchBase: 'Vandenberg Space Force Base',
    status: 'scheduled',
    description: 'Deployment of 40 OneWeb internet satellites.',
  },
  {
    id: 4,
    name: 'Dragon Cargo Mission',
    date: '2025-10-28T22:00:00Z',
    rocket: 'Falcon 9',
    launchBase: 'Kennedy Space Center',
    status: 'scheduled',
    description: 'ISS resupply mission carrying scientific experiments and crew supplies.',
  },
  {
    id: 5,
    name: 'NROL-87',
    date: '2025-11-05T22:45:00Z',
    rocket: 'Falcon Heavy',
    launchBase: 'Kennedy Space Center',
    status: 'scheduled',
    description: 'National Reconnaissance Office payload launch.',
  },
  {
    id: 6,
    name: 'Crew-8 Mission',
    date: '2025-10-20T15:30:00Z',
    rocket: 'Falcon 9',
    launchBase: 'Kennedy Space Center',
    status: 'successful',
    description: 'NASA astronauts and international partners launched to the ISS.',
  },
  {
    id: 7,
    name: 'SES-24 Communications Satellite',
    date: '2025-10-18T03:45:00Z',
    rocket: 'Falcon 9',
    launchBase: 'Cape Canaveral',
    status: 'successful',
    description: 'Commercial communications satellite deployment to geostationary orbit.',
  },
  {
    id: 8,
    name: 'Starship IFT-5',
    date: '2025-10-15T12:00:00Z',
    rocket: 'Starship',
    launchBase: 'Boca Chica Launch Site',
    status: 'successful',
    description: 'Fifth integrated flight test of Starship with booster catch attempt.',
  },
  {
    id: 9,
    name: 'Artemis III',
    date: '2025-12-20T14:00:00Z',
    rocket: 'Starship',
    launchBase: 'Boca Chica Launch Site',
    status: 'scheduled',
    description: 'First crewed lunar landing mission of the Artemis program.',
  },
];

export const news: News[] = [
  {
    id: 1,
    title: 'SpaceX Completes 200th Successful Landing',
    summary: 'SpaceX has achieved another milestone with its 200th successful booster landing, demonstrating the reliability of reusable rocket technology.',
    content: `# SpaceX Achieves Historic 200th Successful Landing

SpaceX has reached a remarkable milestone in aerospace engineering, successfully completing its 200th booster landing. This achievement underscores the company's leadership in developing and perfecting reusable rocket technology.

## A Revolutionary Achievement

The Falcon 9 first-stage booster touched down precisely on the autonomous drone ship "Just Read the Instructions" in the Atlantic Ocean, marking this historic 200th recovery. This landing represents years of innovation, testing, and continuous improvement in rocket reusability.

### Key Statistics

- **Total Landings**: 200 successful recoveries
- **Success Rate**: Over 95% landing success rate
- **Cost Savings**: Estimated billions of dollars saved through reusability
- **Turnaround Time**: Some boosters have flown up to 15 times

## Impact on Space Industry

The success of SpaceX's reusable rocket program has fundamentally changed the economics of space access. By recovering and refurbishing boosters, SpaceX has:

1. Reduced launch costs significantly
2. Increased launch frequency capabilities
3. Demonstrated the viability of reusable space technology
4. Inspired competition and innovation across the industry

> "This milestone represents the culmination of years of hard work by our team. Reusability is key to making life multiplanetary." - SpaceX Engineering Team

## Looking Forward

With 200 successful landings achieved, SpaceX continues to push the boundaries of what's possible. The company is now focusing on:

- Further reducing turnaround time between flights
- Increasing the number of times a single booster can fly
- Applying reusability lessons to the Starship program
- Expanding the fleet of recovery vessels

The 200th landing is not just a number—it's a testament to the dedication of thousands of engineers and the vision of making space more accessible for humanity.`,
    date: '2025-10-20T12:00:00Z',
    url: 'https://www.spacex.com',
    imageUrl: 'https://via.placeholder.com/600x400/1e3a8a/ffffff?text=SpaceX+Landing',
  },
  {
    id: 2,
    title: 'New Launch Pad Construction Begins at Cape Canaveral',
    summary: 'NASA and SpaceX have broken ground on a new launch pad specifically designed for Starship missions.',
    content: `# Construction Begins on New Starship Launch Pad at Cape Canaveral

NASA and SpaceX have officially broken ground on a state-of-the-art launch pad at Cape Canaveral Space Force Station, specifically designed to support Starship missions. This new facility represents a significant investment in the future of deep space exploration.

## Facility Overview

The new Launch Complex will be one of the most advanced launch facilities ever built, featuring:

- **Launch Mount**: Custom-designed for Starship's unique requirements
- **Propellant Storage**: Massive tanks for methane and liquid oxygen
- **Integration Tower**: A 146-meter tall structure for vehicle stacking
- **Flame Trench**: Advanced design to handle Starship's 33 Raptor engines

### Construction Timeline

| Phase | Timeline | Status |
|-------|----------|--------|
| Site Preparation | Q4 2025 | In Progress |
| Foundation Work | Q1 2026 | Scheduled |
| Tower Construction | Q2-Q3 2026 | Scheduled |
| Systems Integration | Q4 2026 | Scheduled |
| Commissioning | Q1 2027 | Planned |

## Strategic Importance

The new launch pad at Cape Canaveral will serve multiple purposes:

1. **NASA Artemis Program**: Support lunar missions
2. **Commercial Launches**: Enable satellite deployment and space station missions
3. **Mars Missions**: Facilitate future Mars exploration efforts
4. **Crew Operations**: Support human spaceflight missions

### Environmental Considerations

The project has undergone extensive environmental review and includes:

- Wildlife protection measures
- Water quality monitoring systems
- Noise reduction technologies
- Sustainable construction practices

## Community Impact

The construction project is expected to:

- Create over 1,500 jobs during construction
- Generate 300 permanent operational positions
- Boost local economy by an estimated $500 million annually
- Enhance Florida's position as a global spaceflight hub

> "This launch pad represents not just infrastructure, but the gateway to humanity's future in space." - NASA Administrator

The new facility is expected to be operational by early 2027, just in time to support an ambitious schedule of lunar and deep space missions.`,
    date: '2025-10-18T14:30:00Z',
    url: 'https://www.nasa.gov',
    imageUrl: 'https://via.placeholder.com/600x400/059669/ffffff?text=Launch+Pad',
  },
  {
    id: 3,
    title: 'Blue Origin Announces New Crew Mission',
    summary: 'Blue Origin has announced a new crew mission for early 2026, marking another step in commercial space tourism.',
    content: `# Blue Origin Announces New Crew Mission for 2026

Blue Origin has announced plans for an exciting new crew mission scheduled for early 2026, continuing the company's commitment to making space accessible to private citizens and researchers.

## Mission Details

The upcoming NS-28 mission will carry six passengers aboard the New Shepard vehicle to the edge of space, providing them with approximately 10 minutes of weightlessness and breathtaking views of Earth.

### Flight Profile

The mission will follow Blue Origin's proven flight profile:

- **Launch**: Vertical takeoff from West Texas facility
- **Ascent**: Approximately 3 minutes to apogee
- **Apogee**: Above 100 km (Kármán line)
- **Microgravity**: 10+ minutes of weightlessness
- **Descent**: Parachute landing in the Texas desert

## Crew Selection

Blue Origin is implementing a unique crew selection process for this mission:

- Two seats reserved for private astronauts
- Two seats for research institutions
- Two seats for Blue Origin's Club for the Future program participants

### Training Program

Selected crew members will undergo comprehensive training including:

1. Mission simulation exercises
2. Safety procedures and emergency protocols
3. Zero-gravity familiarization
4. Vehicle systems overview
5. Physical and psychological preparation

## Scientific Objectives

Beyond space tourism, the mission will conduct several scientific experiments:

- **Materials Science**: Testing new materials in microgravity
- **Biological Research**: Studying cellular behavior in space
- **Earth Observation**: High-resolution photography for climate research
- **Technology Demonstration**: Testing new spacecraft systems

### Research Partnerships

Blue Origin is collaborating with:

- Leading universities for educational outreach
- Research institutions for scientific experiments
- Technology companies for system validation
- Space agencies for data sharing

## Vehicle Enhancements

The New Shepard vehicle for this mission features several upgrades:

- Enhanced crew cabin comfort systems
- Improved viewing windows
- Advanced communication systems
- Upgraded safety features

> "Every mission brings us closer to a future where millions of people are living and working in space." - Blue Origin Team

## Booking and Availability

Interested participants can learn more about future missions through Blue Origin's official website. The company plans to increase flight frequency throughout 2026 and beyond.`,
    date: '2025-10-15T10:00:00Z',
    url: 'https://www.blueorigin.com',
    imageUrl: 'https://via.placeholder.com/600x400/dc2626/ffffff?text=Blue+Origin',
  },
];

export const launchBases: LaunchBase[] = [
  {
    id: 1,
    name: 'Kennedy Space Center',
    location: 'Merritt Island, Florida',
    country: 'United States',
    description: 'NASA\'s primary launch center for human spaceflight, located on the east coast of Florida.',
    imageUrl: 'https://via.placeholder.com/800x600/1e3a8a/ffffff?text=Kennedy+Space+Center',
    latitude: 28.5729,
    longitude: -80.6490,
  },
  {
    id: 2,
    name: 'Cape Canaveral Space Force Station',
    location: 'Cape Canaveral, Florida',
    country: 'United States',
    description: 'A major launch site for military and commercial payloads, adjacent to Kennedy Space Center.',
    imageUrl: 'https://via.placeholder.com/800x600/7c3aed/ffffff?text=Cape+Canaveral',
    latitude: 28.4889,
    longitude: -80.5778,
  },
  {
    id: 3,
    name: 'Boca Chica Launch Site',
    location: 'Boca Chica, Texas',
    country: 'United States',
    description: 'SpaceX\'s private launch facility for Starship development and operations.',
    imageUrl: 'https://via.placeholder.com/800x600/059669/ffffff?text=Boca+Chica',
    latitude: 25.9972,
    longitude: -97.1572,
  },
  {
    id: 4,
    name: 'Vandenberg Space Force Base',
    location: 'Lompoc, California',
    country: 'United States',
    description: 'Primary launch site for polar orbit missions on the west coast.',
    imageUrl: 'https://via.placeholder.com/800x600/ea580c/ffffff?text=Vandenberg',
    latitude: 34.7420,
    longitude: -120.5724,
  },
];

export const companies: Company[] = [
  {
    id: 1,
    name: 'SpaceX',
    description: 'Space Exploration Technologies Corp. is an American aerospace manufacturer and space transportation company founded with the goal of reducing space transportation costs and enabling the colonization of Mars.',
    founded: 2002,
    founder: 'Elon Musk',
    headquarters: 'Hawthorne, California, USA',
    employees: 13000,
    website: 'https://www.spacex.com',
    imageUrl: 'https://via.placeholder.com/400x400/1e3a8a/ffffff?text=SpaceX',
  },
  {
    id: 2,
    name: 'Blue Origin',
    description: 'An American private spaceflight company developing technologies to enable human access to space at lower costs and increased reliability.',
    founded: 2000,
    founder: 'Jeff Bezos',
    headquarters: 'Kent, Washington, USA',
    employees: 10000,
    website: 'https://www.blueorigin.com',
    imageUrl: 'https://via.placeholder.com/400x400/dc2626/ffffff?text=Blue+Origin',
  },
  {
    id: 3,
    name: 'NASA',
    description: 'The National Aeronautics and Space Administration is an independent agency of the U.S. federal government responsible for the civilian space program and aerospace research.',
    founded: 1958,
    founder: 'U.S. Government',
    headquarters: 'Washington, D.C., USA',
    employees: 18000,
    website: 'https://www.nasa.gov',
    imageUrl: 'https://via.placeholder.com/400x400/059669/ffffff?text=NASA',
  },
  {
    id: 4,
    name: 'Rocket Lab',
    description: 'An American aerospace manufacturer with a wholly-owned New Zealand subsidiary, specializing in small satellite launches.',
    founded: 2006,
    founder: 'Peter Beck',
    headquarters: 'Long Beach, California, USA',
    employees: 1700,
    website: 'https://www.rocketlabusa.com',
    imageUrl: 'https://via.placeholder.com/400x400/7c3aed/ffffff?text=Rocket+Lab',
  },
];
