// ── HAMBURGER MENU ──
function toggleMenu() {
    const navLinks = document.querySelector('.nav-links')
    navLinks.classList.toggle('active')
}

// ── CLOSE MENU ON LINK CLICK ──
document.querySelectorAll('.nav-links a').forEach(link => {
    link.addEventListener('click', () => {
        document.querySelector('.nav-links').classList.remove('active')
    })
})

// ── SMOOTH SCROLL ──
document.querySelectorAll('a[href^="#"]').forEach(anchor => {
    anchor.addEventListener('click', function(e) {
        e.preventDefault()
        document.querySelector(this.getAttribute('href')).scrollIntoView({
            behavior: 'smooth'
        })
    })
})

// ── NAVBAR SCROLL EFFECT ──
window.addEventListener('scroll', () => {
    const navbar = document.querySelector('.navbar')
    if (window.scrollY > 50) {
        navbar.style.boxShadow = '0 2px 20px rgba(212,175,55,0.3)'
    } else {
        navbar.style.boxShadow = 'none'
    }
})

// ── FADE IN ON SCROLL ──
const observer = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
        if (entry.isIntersecting) {
            entry.target.classList.add('visible')
        }
    })
}, { threshold: 0.1 })

document.querySelectorAll('.collection-card, .ceo-content').forEach(el => {
    observer.observe(el)
})