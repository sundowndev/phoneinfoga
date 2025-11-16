// Initialize Reveal.js
Reveal.initialize({
    hash: true,
    center: true,
    transition: 'cube', // none/fade/slide/convex/concave/zoom
    transitionSpeed: 'default',
    backgroundTransition: 'zoom',
    slideNumber: true,
    controlsLayout: 'bottom-right',
    controlsBackArrows: 'faded',
    progress: true,
    keyboard: true,
    overview: true,
    touch: true,
    loop: false,
    rtl: true, // Right to left for Arabic
    navigationMode: 'default',
    preloadIframes: true,
    autoAnimate: true,
    autoAnimateDuration: 1.0,
    view: 'scroll',
    scrollActivationWidth: 435,

    // Parallax Background
    parallaxBackgroundImage: '',
    parallaxBackgroundSize: '',
    parallaxBackgroundHorizontal: 200,
    parallaxBackgroundVertical: 50,
});

// Three.js 3D Background
let scene, camera, renderer, particles;

function init3DBackground() {
    // Create scene
    scene = new THREE.Scene();

    // Create camera
    camera = new THREE.PerspectiveCamera(
        75,
        window.innerWidth / window.innerHeight,
        0.1,
        1000
    );
    camera.position.z = 50;

    // Create renderer
    renderer = new THREE.WebGLRenderer({
        alpha: true,
        antialias: true
    });
    renderer.setSize(window.innerWidth, window.innerHeight);
    renderer.setPixelRatio(window.devicePixelRatio);

    // Add renderer to container
    const container = document.getElementById('canvas-container');
    container.appendChild(renderer.domElement);

    // Create particles
    createParticles();

    // Create Instagram logo shapes
    createInstagramShapes();

    // Add lights
    const ambientLight = new THREE.AmbientLight(0xffffff, 0.5);
    scene.add(ambientLight);

    const pointLight = new THREE.PointLight(0xE1306C, 1, 100);
    pointLight.position.set(10, 10, 10);
    scene.add(pointLight);

    // Handle window resize
    window.addEventListener('resize', onWindowResize, false);

    // Start animation
    animate();
}

function createParticles() {
    const geometry = new THREE.BufferGeometry();
    const vertices = [];
    const colors = [];

    // Instagram gradient colors
    const instagramColors = [
        new THREE.Color(0x405DE6),
        new THREE.Color(0x5851DB),
        new THREE.Color(0x833AB4),
        new THREE.Color(0xC13584),
        new THREE.Color(0xE1306C),
        new THREE.Color(0xFD1D1D),
        new THREE.Color(0xF56040),
        new THREE.Color(0xFEDC80)
    ];

    for (let i = 0; i < 2000; i++) {
        const x = (Math.random() - 0.5) * 200;
        const y = (Math.random() - 0.5) * 200;
        const z = (Math.random() - 0.5) * 200;

        vertices.push(x, y, z);

        const color = instagramColors[Math.floor(Math.random() * instagramColors.length)];
        colors.push(color.r, color.g, color.b);
    }

    geometry.setAttribute('position', new THREE.Float32BufferAttribute(vertices, 3));
    geometry.setAttribute('color', new THREE.Float32BufferAttribute(colors, 3));

    const material = new THREE.PointsMaterial({
        size: 2,
        vertexColors: true,
        transparent: true,
        opacity: 0.8,
        blending: THREE.AdditiveBlending
    });

    particles = new THREE.Points(geometry, material);
    scene.add(particles);
}

function createInstagramShapes() {
    // Create floating cubes with Instagram colors
    const colors = [0x405DE6, 0x833AB4, 0xC13584, 0xE1306C, 0xFD1D1D, 0xF56040];

    for (let i = 0; i < 20; i++) {
        const geometry = new THREE.BoxGeometry(3, 3, 3);
        const material = new THREE.MeshPhongMaterial({
            color: colors[Math.floor(Math.random() * colors.length)],
            transparent: true,
            opacity: 0.6,
            shininess: 100
        });

        const cube = new THREE.Mesh(geometry, material);

        cube.position.x = (Math.random() - 0.5) * 100;
        cube.position.y = (Math.random() - 0.5) * 100;
        cube.position.z = (Math.random() - 0.5) * 100;

        cube.rotation.x = Math.random() * Math.PI;
        cube.rotation.y = Math.random() * Math.PI;

        cube.userData.velocity = {
            x: (Math.random() - 0.5) * 0.02,
            y: (Math.random() - 0.5) * 0.02,
            z: (Math.random() - 0.5) * 0.02,
            rx: (Math.random() - 0.5) * 0.02,
            ry: (Math.random() - 0.5) * 0.02
        };

        scene.add(cube);
    }
}

function animate() {
    requestAnimationFrame(animate);

    // Rotate particles
    if (particles) {
        particles.rotation.x += 0.0003;
        particles.rotation.y += 0.0005;
    }

    // Animate cubes
    scene.children.forEach(child => {
        if (child instanceof THREE.Mesh && child.geometry instanceof THREE.BoxGeometry) {
            child.position.x += child.userData.velocity.x;
            child.position.y += child.userData.velocity.y;
            child.position.z += child.userData.velocity.z;

            child.rotation.x += child.userData.velocity.rx;
            child.rotation.y += child.userData.velocity.ry;

            // Boundary check
            if (Math.abs(child.position.x) > 50) child.userData.velocity.x *= -1;
            if (Math.abs(child.position.y) > 50) child.userData.velocity.y *= -1;
            if (Math.abs(child.position.z) > 50) child.userData.velocity.z *= -1;
        }
    });

    renderer.render(scene, camera);
}

function onWindowResize() {
    camera.aspect = window.innerWidth / window.innerHeight;
    camera.updateProjectionMatrix();
    renderer.setSize(window.innerWidth, window.innerHeight);
}

// Mouse interaction
let mouseX = 0;
let mouseY = 0;

document.addEventListener('mousemove', (event) => {
    mouseX = (event.clientX / window.innerWidth) * 2 - 1;
    mouseY = -(event.clientY / window.innerHeight) * 2 + 1;

    if (camera) {
        camera.position.x = mouseX * 5;
        camera.position.y = mouseY * 5;
        camera.lookAt(scene.position);
    }
});

// Initialize 3D background when page loads
window.addEventListener('DOMContentLoaded', init3DBackground);

// Add sound effects (optional)
Reveal.on('slidechanged', event => {
    // You can add sound effects here if needed
    console.log('Slide changed to:', event.indexh);
});

// Add keyboard shortcuts
document.addEventListener('keydown', (event) => {
    // Press 'H' for help
    if (event.key === 'h' || event.key === 'H') {
        alert(`
مفاتيح التحكم:
- الأسهم: التنقل بين الشرائح
- Space: الشريحة التالية
- ESC: نظرة عامة على العرض
- F: وضع ملء الشاشة
- S: عرض ملاحظات المتحدث
        `);
    }
});

// Fullscreen functionality
document.addEventListener('keydown', (event) => {
    if (event.key === 'f' || event.key === 'F') {
        if (!document.fullscreenElement) {
            document.documentElement.requestFullscreen();
        } else {
            if (document.exitFullscreen) {
                document.exitFullscreen();
            }
        }
    }
});

// Add animation on slide entrance
Reveal.on('slidechanged', event => {
    const currentSlide = event.currentSlide;

    // Add entrance animation class
    currentSlide.classList.add('slide-entered');

    // Trigger special effects for specific slides
    if (currentSlide.querySelector('.instagram-icon')) {
        const icon = currentSlide.querySelector('.instagram-icon');
        icon.style.animation = 'none';
        setTimeout(() => {
            icon.style.animation = 'float 3s ease-in-out infinite';
        }, 10);
    }
});

// Performance optimization
if (window.devicePixelRatio > 1) {
    renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
}

// Console welcome message
console.log(`
%c
╔══════════════════════════════════════════╗
║   Instagram 3D Presentation              ║
║   Powered by Reveal.js & Three.js       ║
║   © 2024 - All Rights Reserved          ║
╚══════════════════════════════════════════╝
`, 'color: #E1306C; font-weight: bold; font-size: 12px;');
